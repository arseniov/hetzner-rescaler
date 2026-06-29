package store

import (
	"testing"
)

func TestServerCreateAndGet(t *testing.T) {
	s := newTestStore(t)
	p, _ := s.CreateProject("p", []byte("t"), []byte("n"))

	srv, err := s.CreateServer(p.ID, Server{
		HCloudServerID: 123,
		Name:           "web",
		BaseServerType: "cpx11",
		TopServerType:  "cpx21",
		FallbackChain:  []string{"cpx21", "cpx11"},
		Mode:           "scheduled",
		Timezone:       "Europe/Rome",
	})
	if err != nil {
		t.Fatalf("CreateServer: %v", err)
	}
	if srv.ID == 0 || srv.HCloudServerID != 123 {
		t.Fatalf("got %+v", srv)
	}

	got, err := s.GetServer(srv.ID)
	if err != nil {
		t.Fatalf("GetServer: %v", err)
	}
	if got.Name != "web" || got.Mode != "scheduled" {
		t.Fatalf("got %+v", got)
	}
	if got.Timezone != "Europe/Rome" {
		t.Fatalf("timezone = %q, want Europe/Rome", got.Timezone)
	}
}

func TestServerListByProject(t *testing.T) {
	s := newTestStore(t)
	p1, _ := s.CreateProject("p1", []byte("t"), []byte("n"))
	p2, _ := s.CreateProject("p2", []byte("t"), []byte("n"))
	_, _ = s.CreateServer(p1.ID, Server{HCloudServerID: 1, Name: "a", BaseServerType: "cpx11", TopServerType: "cpx21", FallbackChain: []string{"cpx21", "cpx11"}, Mode: "manual", Timezone: "UTC"})
	_, _ = s.CreateServer(p1.ID, Server{HCloudServerID: 2, Name: "b", BaseServerType: "cpx11", TopServerType: "cpx21", FallbackChain: []string{"cpx21", "cpx11"}, Mode: "manual", Timezone: "UTC"})
	_, _ = s.CreateServer(p2.ID, Server{HCloudServerID: 3, Name: "c", BaseServerType: "cpx11", TopServerType: "cpx21", FallbackChain: []string{"cpx21", "cpx11"}, Mode: "manual", Timezone: "UTC"})

	got, err := s.ListServersByProject(p1.ID)
	if err != nil {
		t.Fatalf("ListServersByProject: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("got %d, want 2", len(got))
	}
}

func TestWindowCRUD(t *testing.T) {
	s := newTestStore(t)
	p, _ := s.CreateProject("p", []byte("t"), []byte("n"))
	srv, _ := s.CreateServer(p.ID, Server{HCloudServerID: 1, Name: "x", BaseServerType: "cpx11", TopServerType: "cpx21", FallbackChain: []string{"cpx21", "cpx11"}, Mode: "scheduled", Timezone: "UTC"})

	w, err := s.CreateWindow(srv.ID, Window{
		Label:       "weekday 9-19",
		DaysOfWeek:  0b0011111,
		StartTime:   "09:00",
		StopTime:    "19:00",
		TargetType:  "cpx21",
		Enabled:     true,
	})
	if err != nil {
		t.Fatalf("CreateWindow: %v", err)
	}
	if w.ID == 0 {
		t.Fatal("window id not set")
	}

	wins, err := s.ListWindows(srv.ID)
	if err != nil {
		t.Fatalf("ListWindows: %v", err)
	}
	if len(wins) != 1 || wins[0].Label != "weekday 9-19" {
		t.Fatalf("got %+v", wins)
	}

	if err := s.DeleteWindow(w.ID); err != nil {
		t.Fatalf("DeleteWindow: %v", err)
	}
	wins, _ = s.ListWindows(srv.ID)
	if len(wins) != 0 {
		t.Fatalf("after delete: got %d, want 0", len(wins))
	}
}

func TestWindowLabelUniquePerServer(t *testing.T) {
	s := newTestStore(t)
	p, _ := s.CreateProject("p", []byte("t"), []byte("n"))
	srv, _ := s.CreateServer(p.ID, Server{HCloudServerID: 1, Name: "x", BaseServerType: "cpx11", TopServerType: "cpx21", FallbackChain: []string{"cpx21", "cpx11"}, Mode: "scheduled", Timezone: "UTC"})

	w := Window{Label: "morning", DaysOfWeek: 127, StartTime: "09:00", StopTime: "19:00", TargetType: "cpx21", Enabled: true}
	if _, err := s.CreateWindow(srv.ID, w); err != nil {
		t.Fatalf("first: %v", err)
	}
	if _, err := s.CreateWindow(srv.ID, w); err == nil {
		t.Fatal("expected duplicate label to fail")
	}
}

func TestServerUniquePerProject(t *testing.T) {
	s := newTestStore(t)
	p, _ := s.CreateProject("p", []byte("t"), []byte("n"))
	srv := Server{HCloudServerID: 1, Name: "x", BaseServerType: "cpx11", TopServerType: "cpx21", FallbackChain: []string{"cpx21", "cpx11"}, Mode: "manual", Timezone: "UTC"}
	if _, err := s.CreateServer(p.ID, srv); err != nil {
		t.Fatalf("first: %v", err)
	}
	if _, err := s.CreateServer(p.ID, srv); err == nil {
		t.Fatal("expected duplicate (project, hcloud_server_id) to fail")
	}
}