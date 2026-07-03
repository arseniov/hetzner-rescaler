package api

import (
	"net/http"
	"strconv"

	"github.com/jonamat/hetzner-rescaler/internal/store"
)

func (d Deps) handleServerEvents(w http.ResponseWriter, r *http.Request) {
	sid, ok := pathInt64(r, "id")
	if !ok {
		writeJSONError(w, http.StatusBadRequest, "invalid id")
		return
	}
	limit := parseLimit(r)
	events, err := d.Store.ListEventsByServer(sid, limit)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	out := make([]EventResponse, 0, len(events))
	for _, e := range events {
		out = append(out, eventToResponse(e))
	}
	writeJSON(w, http.StatusOK, out)
}

func (d Deps) handleGlobalEvents(w http.ResponseWriter, r *http.Request) {
	var serverID *int64
	if raw := r.URL.Query().Get("server_id"); raw != "" {
		v, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid server_id")
			return
		}
		serverID = &v
	}
	limit := parseLimit(r)
	events, err := d.Store.ListAllEvents(limit, serverID)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	out := make([]EventResponse, 0, len(events))
	for _, e := range events {
		out = append(out, eventToResponse(e))
	}
	writeJSON(w, http.StatusOK, out)
}

// parseLimit reads the optional "limit" query parameter. Defaults to 100
// when missing or unparseable to bound the response size on public endpoints.
func parseLimit(r *http.Request) int {
	if raw := r.URL.Query().Get("limit"); raw != "" {
		v, err := strconv.Atoi(raw)
		if err == nil && v > 0 {
			return v
		}
	}
	return 100
}

func eventToResponse(e *store.Event) EventResponse {
	return EventResponse{
		ID:          e.ID,
		ServerID:    e.ServerID,
		Kind:        e.Kind,
		FromType:    e.FromType,
		ToType:      e.ToType,
		StartedAt:   e.StartedAt,
		FinishedAt:  e.FinishedAt,
		OK:          e.OK,
		Error:       e.Error,
		TriggeredBy: e.TriggeredBy,
	}
}
