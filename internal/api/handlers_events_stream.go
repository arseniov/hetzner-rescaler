package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

// eventsStreamAuth allows EventSource (which cannot set custom headers)
// to deliver the shared secret via ?token=… in addition to the standard
// X-Internal-Token header. If the header is absent but a non-empty
// `token` query parameter is present, the parameter is copied into the
// header before delegating to the normal internal-token middleware.
//
// Use this wrapper only for routes consumed by EventSource. All other
// /api/* routes continue to require the header via auth().
func eventsStreamAuth(token string, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Internal-Token") == "" {
			if q := r.URL.Query().Get("token"); q != "" {
				r.Header.Set("X-Internal-Token", q)
			}
		}
		RequireInternalToken(token)(h).ServeHTTP(w, r)
	})
}

// handleEventsStream streams live events to the client using Server-Sent
// Events. The connection holds open until the client disconnects.
//
// Protocol:
//   - On connect: an immediate "ready" event with body {"ok":true}.
//   - Each broadcast event is delivered as an "event: rescale" with the
//     JSON-encoded EventResponse as its data payload.
//   - A ": keepalive" comment is sent every 25 seconds to keep proxies
//     and browsers from idling out the connection.
//
// The handler subscribes to the store's broadcast hub with buffer 32; if
// the hub is not attached the handler still emits the ready event and
// keeps the connection open with keepalives (graceful degradation).
func (d Deps) handleEventsStream(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		writeJSONError(w, http.StatusInternalServerError, "streaming unsupported")
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")
	w.WriteHeader(http.StatusOK)

	// Write the ready event immediately and flush so the client knows the
	// stream is live even if no events arrive for a while.
	if _, err := w.Write([]byte("event: ready\ndata: {\"ok\":true}\n\n")); err != nil {
		return
	}
	flusher.Flush()

	hub := d.Store.EventHub()
	if hub == nil {
		// No hub attached — just keep the connection open with keepalives
		// so the client can see the stream is healthy. This is the
		// graceful-degradation path when the broadcaster was never wired
		// (e.g. during certain test setups).
		keepAliveUntilDone(r.Context(), w, flusher)
		return
	}

	ch, unsub := hub.Subscribe(32)
	defer unsub()

	keepAlive := time.NewTicker(25 * time.Second)
	defer keepAlive.Stop()

	for {
		select {
		case <-r.Context().Done():
			return
		case ev, open := <-ch:
			if !open {
				return
			}
			payload, err := json.Marshal(eventToResponse(&ev))
			if err != nil {
				// Should not happen for EventResponse, but skip the event
				// rather than tearing down the stream on a marshal error.
				continue
			}
			if _, err := w.Write([]byte("event: rescale\ndata: ")); err != nil {
				return
			}
			if _, err := w.Write(payload); err != nil {
				return
			}
			if _, err := w.Write([]byte("\n\n")); err != nil {
				return
			}
			flusher.Flush()
		case <-keepAlive.C:
			if _, err := w.Write([]byte(": keepalive\n\n")); err != nil {
				return
			}
			flusher.Flush()
		}
	}
}

// keepAliveUntilDone writes an SSE comment keepalive every 25 seconds
// until the request context is cancelled. Used when the store has no
// broadcast hub attached.
func keepAliveUntilDone(ctx context.Context, w http.ResponseWriter, flusher http.Flusher) {
	keepAlive := time.NewTicker(25 * time.Second)
	defer keepAlive.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-keepAlive.C:
			if _, err := w.Write([]byte(": keepalive\n\n")); err != nil {
				return
			}
			flusher.Flush()
		}
	}
}