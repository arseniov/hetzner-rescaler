package store

import (
	"path/filepath"
	"testing"
)

func newTestStore(t *testing.T) *Store {
	t.Helper()
	s, err := Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	t.Cleanup(func() { s.Close() })
	return s
}

func TestProjectCreateAndGet(t *testing.T) {
	s := newTestStore(t)
	p, err := s.CreateProject("alpha", []byte("tok-enc"), []byte("nonce12byts"))
	if err != nil {
		t.Fatalf("CreateProject: %v", err)
	}
	if p.ID == 0 || p.Name != "alpha" {
		t.Fatalf("got %+v, want ID>0 and name=alpha", p)
	}

	got, err := s.GetProject(p.ID)
	if err != nil {
		t.Fatalf("GetProject: %v", err)
	}
	if got.Name != "alpha" {
		t.Fatalf("got name=%q, want alpha", got.Name)
	}
	if string(got.HCloudTokenEncrypted) != "tok-enc" {
		t.Fatalf("token not round-tripped")
	}
}

func TestProjectList(t *testing.T) {
	s := newTestStore(t)
	_, _ = s.CreateProject("a", []byte("t1"), []byte("n1"))
	_, _ = s.CreateProject("b", []byte("t2"), []byte("n2"))

	projects, err := s.ListProjects()
	if err != nil {
		t.Fatalf("ListProjects: %v", err)
	}
	if len(projects) != 2 {
		t.Fatalf("got %d projects, want 2", len(projects))
	}
}

func TestProjectDeleteCascadesServers(t *testing.T) {
	s := newTestStore(t)
	p, _ := s.CreateProject("p", []byte("t"), []byte("n"))
	srv, err := s.CreateServer(p.ID, Server{
		HCloudServerID:  123,
		Name:            "web",
		BaseServerType:  "cpx11",
		TopServerType:   "cpx21",
		FallbackChain:   []string{"cpx21", "cpx11"},
		Mode:            "scheduled",
		Timezone:        "UTC",
	})
	if err != nil {
		t.Fatalf("CreateServer: %v", err)
	}
	if err := s.DeleteProject(p.ID); err != nil {
		t.Fatalf("DeleteProject: %v", err)
	}
	if _, err := s.GetServer(srv.ID); err == nil {
		t.Fatal("expected GetServer to fail after project delete (cascade)")
	}
}

func TestProjectNameUniqueness(t *testing.T) {
	s := newTestStore(t)
	if _, err := s.CreateProject("dup", []byte("t"), []byte("n")); err != nil {
		t.Fatalf("first CreateProject: %v", err)
	}
	if _, err := s.CreateProject("dup", []byte("t2"), []byte("n2")); err == nil {
		t.Fatal("expected error on duplicate name")
	}
}