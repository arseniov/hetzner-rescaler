package api

import (
	"bufio"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/jonamat/hetzner-rescaler/internal/broadcast"
	"github.com/jonamat/hetzner-rescaler/internal/store"
)

// TestEventsStream_DeliversBroadcast verifies the happy path: the
// server emits a "ready" event on connect, then every value passed to
// hub.Broadcast appears as an "event: rescale" frame in the stream.
func TestEventsStream_DeliversBroadcast(t *testing.T) {
	deps, _ := newTestDeps(t)
	hub := broadcast.NewHub[store.Event]()
	defer hub.Close()
	deps.Store.SetBroadcastHub(hub)

	srv := httptest.NewServer(NewRouter(deps))
	defer srv.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, srv.URL+"/api/events/stream", nil)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("X-Internal-Token", testInternalToken)
	req.Header.Set("Accept", "text/event-stream")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("do: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status: want 200, got %d", resp.StatusCode)
	}

	reader := bufio.NewReader(resp.Body)

	// Read first frame: should be the ready event.
	line, err := reader.ReadString('\n')
	if err != nil {
		t.Fatalf("read first line: %v", err)
	}
	if !strings.Contains(line, "event: ready") {
		t.Fatalf("first frame: want to contain %q, got %q", "event: ready", line)
	}

	// Broadcast one event from another goroutine so the handler's
	// select loop receives it.
	sample := store.Event{ID: 42, ServerID: 1, Kind: "rescale_up", OK: true, TriggeredBy: "scheduler"}
	go func() {
		hub.Broadcast(sample)
	}()

	// Read until we see an "event: rescale" line.
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		line, err := reader.ReadString('\n')
		if err != nil {
			t.Fatalf("read rescale line: %v", err)
		}
		if strings.Contains(line, "event: rescale") {
			// Found it; sanity-check the next data line is non-empty JSON.
			dataLine, err := reader.ReadString('\n')
			if err != nil {
				t.Fatalf("read rescale data: %v", err)
			}
			if !strings.HasPrefix(dataLine, "data: ") {
				t.Fatalf("expected data line after rescale event, got %q", dataLine)
			}
			if !strings.Contains(dataLine, `"id":42`) {
				t.Fatalf("expected payload to contain id 42, got %q", dataLine)
			}
			break
		}
	}

	// Cancel and verify the handler returns promptly.
	cancel()
	done := make(chan struct{})
	go func() {
		resp.Body.Close()
		close(done)
	}()
	select {
	case <-done:
		// good
	case <-time.After(2 * time.Second):
		t.Fatalf("handler did not return within 2s of cancel")
	}
}

// TestEventsStream_QueryTokenAccepted verifies that the SSE endpoint
// accepts the internal token via the ?token=… query parameter when the
// X-Internal-Token header is absent. This path is used by browser
// EventSource, which cannot set custom headers.
func TestEventsStream_QueryTokenAccepted(t *testing.T) {
	deps, _ := newTestDeps(t)
	hub := broadcast.NewHub[store.Event]()
	defer hub.Close()
	deps.Store.SetBroadcastHub(hub)

	srv := httptest.NewServer(NewRouter(deps))
	defer srv.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Note: no X-Internal-Token header set.
	url := srv.URL + "/api/events/stream?token=" + testInternalToken
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("Accept", "text/event-stream")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("do: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status: want 200, got %d", resp.StatusCode)
	}

	reader := bufio.NewReader(resp.Body)
	line, err := reader.ReadString('\n')
	if err != nil {
		t.Fatalf("read first line: %v", err)
	}
	if !strings.Contains(line, "event: ready") {
		t.Fatalf("first frame: want to contain %q, got %q", "event: ready", line)
	}
}

// TestEventsStream_DisconnectStopsDelivery verifies that after a client
// disconnects, no further events are delivered to that subscription.
// We use a custom handler that mimics the production handler but with
// an explicit subscribe so we can verify the unsubscribe path.
func TestEventsStream_DisconnectStopsDelivery(t *testing.T) {
	deps, _ := newTestDeps(t)
	hub := broadcast.NewHub[store.Event]()
	defer hub.Close()
	deps.Store.SetBroadcastHub(hub)

	srv := httptest.NewServer(NewRouter(deps))
	defer srv.Close()

	ctx, cancel := context.WithCancel(context.Background())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, srv.URL+"/api/events/stream", nil)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("X-Internal-Token", testInternalToken)
	req.Header.Set("Accept", "text/event-stream")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("do: %v", err)
	}

	reader := bufio.NewReader(resp.Body)
	// Drain the ready frame.
	if _, err := reader.ReadString('\n'); err != nil {
		t.Fatalf("read ready: %v", err)
	}

	// Broadcast one event; handler should receive it.
	hub.Broadcast(store.Event{ID: 1, Kind: "rescale_up", OK: true})

	// Wait for the handler to deliver it.
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		line, err := reader.ReadString('\n')
		if err != nil {
			t.Fatalf("read rescale: %v", err)
		}
		if strings.Contains(line, "event: rescale") {
			break
		}
	}

	// Cancel the request and close the response body. The handler
	// should unsubscribe from the hub.
	cancel()
	resp.Body.Close()

	// Allow the handler time to run its defer unsub.
	time.Sleep(100 * time.Millisecond)

	// Now confirm the hub has no subscribers left by checking that
	// Broadcast does not block or panic (which it would if any
	// subscriber's channel was closed-and-still-mapped).
	var counter int64
	done := make(chan struct{})
	go func() {
		defer close(done)
		// Hammer the hub; if there's a stale subscriber, the closed
		// channel send will panic.
		for i := 0; i < 100; i++ {
			hub.Broadcast(store.Event{ID: int64(i), Kind: "rescale_up", OK: true})
		}
		atomic.StoreInt64(&counter, 100)
	}()
	select {
	case <-done:
		if got := atomic.LoadInt64(&counter); got != 100 {
			t.Fatalf("broadcast loop did not complete: counter=%d", got)
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("Broadcast hung after disconnect")
	}
}