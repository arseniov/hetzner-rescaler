package store

import (
	"testing"
	"time"
)

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
	got, err := s.ListAllEvents(0, &id1, 100)
	if err != nil {
		t.Fatalf("ListAllEvents: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("got %d, want 1", len(got))
	}

	// No filter
	got, err = s.ListAllEvents(0, nil, 100)
	if err != nil {
		t.Fatalf("ListAllEvents: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("got %d, want 2", len(got))
	}
}