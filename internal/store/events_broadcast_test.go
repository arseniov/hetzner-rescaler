package store

import (
	"testing"
	"time"

	"github.com/jonamat/hetzner-rescaler/internal/broadcast"
)

// TestAppendEventBroadcastsToHub verifies that a subscriber to the broadcast
// hub receives every event appended via Store.AppendEvent, with the row's
// assigned ID populated on the delivered value.
func TestAppendEventBroadcastsToHub(t *testing.T) {
	s := newTestStore(t)
	p, _ := s.CreateProject("p", []byte("t"), []byte("n"))
	srv, _ := s.CreateServer(p.ID, Server{HCloudServerID: 1, Name: "x", BaseServerType: "cpx11", TopServerType: "cpx21", FallbackChain: []string{"cpx21", "cpx11"}, Mode: "manual", Timezone: "UTC"})

	hub := broadcast.NewHub[Event]()
	ch, unsub := hub.Subscribe(4)
	defer unsub()
	s.SetBroadcastHub(hub)

	startedAt := time.Now().UTC()
	finishedAt := startedAt.Add(time.Second)
	id, err := s.AppendEvent(Event{
		ServerID:    srv.ID,
		Kind:        "rescale_up",
		FromType:    "cpx11",
		ToType:      "cpx21",
		StartedAt:   startedAt,
		FinishedAt:  finishedAt,
		OK:          true,
		TriggeredBy: "scheduler",
	})
	if err != nil {
		t.Fatalf("AppendEvent: %v", err)
	}
	if id == 0 {
		t.Fatal("event id not set")
	}

	select {
	case got := <-ch:
		if got.ID != id {
			t.Fatalf("delivered event ID = %d, want %d", got.ID, id)
		}
		if got.ServerID != srv.ID {
			t.Fatalf("delivered event ServerID = %d, want %d", got.ServerID, srv.ID)
		}
		if got.Kind != "rescale_up" || got.FromType != "cpx11" || got.ToType != "cpx21" {
			t.Fatalf("delivered event fields = %+v", got)
		}
		if got.TriggeredBy != "scheduler" {
			t.Fatalf("delivered event TriggeredBy = %q, want %q", got.TriggeredBy, "scheduler")
		}
		if !got.OK {
			t.Fatal("delivered event OK = false, want true")
		}
		if !got.StartedAt.Equal(startedAt) {
			t.Fatalf("delivered event StartedAt = %v, want %v", got.StartedAt, startedAt)
		}
		if !got.FinishedAt.Equal(finishedAt) {
			t.Fatalf("delivered event FinishedAt = %v, want %v", got.FinishedAt, finishedAt)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for broadcast")
	}
}

// TestAppendEventNoHubIsNoop verifies that AppendEvent works fine when no
// hub has been attached (the default for tests / CLI-only flows).
func TestAppendEventNoHubIsNoop(t *testing.T) {
	s := newTestStore(t)
	p, _ := s.CreateProject("p", []byte("t"), []byte("n"))
	srv, _ := s.CreateServer(p.ID, Server{HCloudServerID: 1, Name: "x", BaseServerType: "cpx11", TopServerType: "cpx21", FallbackChain: []string{"cpx21", "cpx11"}, Mode: "manual", Timezone: "UTC"})

	// No SetBroadcastHub call: s.hub is nil.
	if _, err := s.AppendEvent(Event{ServerID: srv.ID, Kind: "noop", StartedAt: time.Now().UTC(), OK: true, TriggeredBy: "test"}); err != nil {
		t.Fatalf("AppendEvent: %v", err)
	}
}

// TestAppendEventUnsubscribeStopsDelivery verifies that unsubscribing a
// channel stops further deliveries.
func TestAppendEventUnsubscribeStopsDelivery(t *testing.T) {
	s := newTestStore(t)
	p, _ := s.CreateProject("p", []byte("t"), []byte("n"))
	srv, _ := s.CreateServer(p.ID, Server{HCloudServerID: 1, Name: "x", BaseServerType: "cpx11", TopServerType: "cpx21", FallbackChain: []string{"cpx21", "cpx11"}, Mode: "manual", Timezone: "UTC"})

	hub := broadcast.NewHub[Event]()
	ch, unsub := hub.Subscribe(4)
	s.SetBroadcastHub(hub)

	// First event should arrive.
	if _, err := s.AppendEvent(Event{ServerID: srv.ID, Kind: "a", StartedAt: time.Now().UTC(), OK: true, TriggeredBy: "test"}); err != nil {
		t.Fatalf("AppendEvent: %v", err)
	}
	select {
	case <-ch:
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for first broadcast")
	}

	// Unsubscribe, then insert another event: it must NOT be delivered.
	unsub()
	if _, err := s.AppendEvent(Event{ServerID: srv.ID, Kind: "b", StartedAt: time.Now().UTC(), OK: true, TriggeredBy: "test"}); err != nil {
		t.Fatalf("AppendEvent: %v", err)
	}
	select {
	case ev, ok := <-ch:
		if ok {
			t.Fatalf("got delivery after unsubscribe: %+v", ev)
		}
		// Channel closed by unsub; that's the expected signal.
	case <-time.After(200 * time.Millisecond):
		// No event and channel still open (or closed without delivery) — both fine.
		// The strong guarantee is no delivery: we already see it.
	}
}

// TestAppendEventDetachingHubStopsDelivery verifies that passing nil to
// SetBroadcastHub detaches the hub.
func TestAppendEventDetachingHubStopsDelivery(t *testing.T) {
	s := newTestStore(t)
	p, _ := s.CreateProject("p", []byte("t"), []byte("n"))
	srv, _ := s.CreateServer(p.ID, Server{HCloudServerID: 1, Name: "x", BaseServerType: "cpx11", TopServerType: "cpx21", FallbackChain: []string{"cpx21", "cpx11"}, Mode: "manual", Timezone: "UTC"})

	hub := broadcast.NewHub[Event]()
	ch, unsub := hub.Subscribe(4)
	defer unsub()
	s.SetBroadcastHub(hub)

	if _, err := s.AppendEvent(Event{ServerID: srv.ID, Kind: "a", StartedAt: time.Now().UTC(), OK: true, TriggeredBy: "test"}); err != nil {
		t.Fatalf("AppendEvent: %v", err)
	}
	select {
	case <-ch:
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for first broadcast")
	}

	s.SetBroadcastHub(nil)
	if _, err := s.AppendEvent(Event{ServerID: srv.ID, Kind: "b", StartedAt: time.Now().UTC(), OK: true, TriggeredBy: "test"}); err != nil {
		t.Fatalf("AppendEvent: %v", err)
	}
	select {
	case ev := <-ch:
		t.Fatalf("got delivery after detach: %+v", ev)
	case <-time.After(200 * time.Millisecond):
	}
}