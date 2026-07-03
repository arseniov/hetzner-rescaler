package api

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/jonamat/hetzner-rescaler/internal/store"
)

func TestServerEvents_ReturnsArray(t *testing.T) {
	deps, _ := newTestDeps(t)
	h := NewRouter(deps)
	_, sid := seedServer(t, deps, "p1", "web-1")
	// Seed three events for the server.
	for _, kind := range []string{"rescale_up", "rescale_down", "rescale_failed"} {
		_, err := deps.Store.AppendEvent(store.Event{
			ServerID: sid, Kind: kind, OK: true, TriggeredBy: "scheduler",
		})
		if err != nil {
			t.Fatalf("AppendEvent: %v", err)
		}
	}

	req := authedRequest(t, "GET", "/api/servers/"+itoa(sid)+"/events", nil)
	rr := recorder(t, h, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", rr.Code)
	}
	var got []EventResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(got) != 3 {
		t.Fatalf("want 3 events, got %d", len(got))
	}
}

func TestServerEvents_RespectsLimit(t *testing.T) {
	deps, _ := newTestDeps(t)
	h := NewRouter(deps)
	_, sid := seedServer(t, deps, "p1", "web-1")
	for i := 0; i < 5; i++ {
		_, _ = deps.Store.AppendEvent(store.Event{
			ServerID: sid, Kind: "rescale_up", OK: true, TriggeredBy: "scheduler",
		})
	}
	req := authedRequest(t, "GET", "/api/servers/"+itoa(sid)+"/events?limit=2", nil)
	rr := recorder(t, h, req)
	var got []EventResponse
	_ = json.Unmarshal(rr.Body.Bytes(), &got)
	if len(got) != 2 {
		t.Fatalf("want 2 events with limit=2, got %d", len(got))
	}
}

func TestGlobalEvents_FiltersByServerID(t *testing.T) {
	deps, _ := newTestDeps(t)
	h := NewRouter(deps)
	_, s1 := seedServer(t, deps, "p1", "web-1")
	_, s2 := seedServer(t, deps, "p2", "web-2")
	for _, sid := range []int64{s1, s2} {
		_, _ = deps.Store.AppendEvent(store.Event{
			ServerID: sid, Kind: "rescale_up", OK: true, TriggeredBy: "scheduler",
		})
	}
	req := authedRequest(t, "GET", "/api/events?server_id="+itoa(s1), nil)
	rr := recorder(t, h, req)
	var got []EventResponse
	_ = json.Unmarshal(rr.Body.Bytes(), &got)
	if len(got) != 1 || got[0].ServerID != s1 {
		t.Fatalf("want only server 1's event, got %+v", got)
	}
}
