package store

import (
	"testing"
	"time"
)

func seedEvent(t *testing.T, s *Store, serverID int64, kind string) int64 {
	t.Helper()
	id, err := s.AppendEvent(Event{
		ServerID:    serverID,
		Kind:        kind,
		StartedAt:   time.Now().UTC(),
		TriggeredBy: "test",
	})
	if err != nil {
		t.Fatalf("AppendEvent: %v", err)
	}
	return id
}

func seedProjectAndServer(t *testing.T, s *Store) (int64, int64) {
	t.Helper()
	p, err := s.CreateProject("p", []byte("tok"), []byte("nonce12byts"))
	if err != nil {
		t.Fatalf("CreateProject: %v", err)
	}
	srv, err := s.CreateServer(p.ID, Server{
		HCloudServerID: 1, Name: "w", Label: "w",
		BaseServerType: "cpx11", TopServerType: "cpx31",
		FallbackChain: []string{"cpx31", "cpx11"},
		Mode: "manual", Timezone: "UTC",
	})
	if err != nil {
		t.Fatalf("CreateServer: %v", err)
	}
	return p.ID, srv.ID
}

func TestUpdateEventPhase_SetsColumn(t *testing.T) {
	s, err := OpenTemp()
	if err != nil {
		t.Fatalf("OpenTemp: %v", err)
	}
	defer s.Close()
	_, srvID := seedProjectAndServer(t, s)
	id := seedEvent(t, s, srvID, "rescale_pending")

	if err := s.UpdateEventPhase(id, "shutting_down"); err != nil {
		t.Fatalf("UpdateEventPhase: %v", err)
	}

	events, err := s.ListEventsByServer(srvID, 10)
	if err != nil {
		t.Fatalf("ListEventsByServer: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("want 1 event, got %d", len(events))
	}
	if events[0].Phase != "shutting_down" {
		t.Fatalf("Phase = %q, want shutting_down", events[0].Phase)
	}
	if !events[0].FinishedAt.IsZero() {
		t.Fatalf("FinishedAt should still be zero, got %v", events[0].FinishedAt)
	}
}

func TestEventAppendAndList(t *testing.T) {
	s := newTestStore(t)
	p, _ := s.CreateProject("p", []byte("t"), []byte("n"))
	srv, _ := s.CreateServer(p.ID, Server{HCloudServerID: 1, Name: "x", BaseServerType: "cpx11", TopServerType: "cpx21", FallbackChain: []string{"cpx21", "cpx11"}, Mode: "manual", Timezone: "UTC"})

	id, err := s.AppendEvent(Event{
		ServerID:    srv.ID,
		Kind:        "rescale_up",
		FromType:    "cpx11",
		ToType:      "cpx21",
		StartedAt:   time.Now().UTC(),
		FinishedAt:  time.Now().UTC(),
		OK:          true,
		TriggeredBy: "scheduler",
	})
	if err != nil {
		t.Fatalf("AppendEvent: %v", err)
	}
	if id == 0 {
		t.Fatal("event id not set")
	}

	events, err := s.ListEventsByServer(srv.ID, 10)
	if err != nil {
		t.Fatalf("ListEventsByServer: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("got %d, want 1", len(events))
	}
	if events[0].Kind != "rescale_up" || events[0].FromType != "cpx11" || events[0].ToType != "cpx21" {
		t.Fatalf("got %+v", events[0])
	}
}

func TestEventListRespectsLimit(t *testing.T) {
	s := newTestStore(t)
	p, _ := s.CreateProject("p", []byte("t"), []byte("n"))
	srv, _ := s.CreateServer(p.ID, Server{HCloudServerID: 1, Name: "x", BaseServerType: "cpx11", TopServerType: "cpx21", FallbackChain: []string{"cpx21", "cpx11"}, Mode: "manual", Timezone: "UTC"})

	for i := 0; i < 5; i++ {
		_, err := s.AppendEvent(Event{ServerID: srv.ID, Kind: "noop", StartedAt: time.Now().UTC(), OK: true, TriggeredBy: "test"})
		if err != nil {
			t.Fatalf("AppendEvent: %v", err)
		}
	}
	got, err := s.ListEventsByServer(srv.ID, 3)
	if err != nil {
		t.Fatalf("ListEventsByServer: %v", err)
	}
	if len(got) != 3 {
		t.Fatalf("got %d, want 3 (limit)", len(got))
	}
}

func TestListAllEventsFilterByServerID(t *testing.T) {
	s := newTestStore(t)
	p, _ := s.CreateProject("p", []byte("t"), []byte("n"))
	s1, _ := s.CreateServer(p.ID, Server{HCloudServerID: 1, Name: "a", BaseServerType: "cpx11", TopServerType: "cpx21", FallbackChain: []string{"cpx21", "cpx11"}, Mode: "manual", Timezone: "UTC"})
	s2, _ := s.CreateServer(p.ID, Server{HCloudServerID: 2, Name: "b", BaseServerType: "cpx11", TopServerType: "cpx21", FallbackChain: []string{"cpx21", "cpx11"}, Mode: "manual", Timezone: "UTC"})

	_, _ = s.AppendEvent(Event{ServerID: s1.ID, Kind: "x", StartedAt: time.Now().UTC(), OK: true, TriggeredBy: "test"})
	_, _ = s.AppendEvent(Event{ServerID: s2.ID, Kind: "y", StartedAt: time.Now().UTC(), OK: true, TriggeredBy: "test"})

	// Filter for s1
	id1 := s1.ID
	got, err := s.ListAllEvents(0, &id1)
	if err != nil {
		t.Fatalf("ListAllEvents: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("got %d, want 1", len(got))
	}

	// No filter
	got, err = s.ListAllEvents(0, nil)
	if err != nil {
		t.Fatalf("ListAllEvents: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("got %d, want 2", len(got))
	}
}

func TestEventListPendingFinishedAtZero(t *testing.T) {
	s := newTestStore(t)
	p, _ := s.CreateProject("p", []byte("t"), []byte("n"))
	srv, _ := s.CreateServer(p.ID, Server{HCloudServerID: 1, Name: "x", BaseServerType: "cpx11", TopServerType: "cpx21", FallbackChain: []string{"cpx21", "cpx11"}, Mode: "manual", Timezone: "UTC"})

	// Append an event without setting FinishedAt — it is still pending.
	_, err := s.AppendEvent(Event{
		ServerID:    srv.ID,
		Kind:        "rescale_up",
		FromType:    "cpx11",
		ToType:      "cpx21",
		StartedAt:   time.Now().UTC(),
		OK:          true,
		TriggeredBy: "scheduler",
	})
	if err != nil {
		t.Fatalf("AppendEvent: %v", err)
	}

	events, err := s.ListEventsByServer(srv.ID, 10)
	if err != nil {
		t.Fatalf("ListEventsByServer: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("got %d, want 1", len(events))
	}
	if !events[0].FinishedAt.IsZero() {
		t.Fatalf("FinishedAt = %v, want zero (pending)", events[0].FinishedAt)
	}
}
