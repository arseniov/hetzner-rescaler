package store

import (
	"errors"
	"testing"
	"time"
)

func newServerForAction(t *testing.T) *Server {
	t.Helper()
	s := newTestStore(t)
	p, _ := s.CreateProject("p", []byte("t"), []byte("n"))
	srv, err := s.CreateServer(p.ID, Server{
		HCloudServerID: 1, Name: "x", BaseServerType: "cpx11", TopServerType: "cpx21",
		FallbackChain: []string{"cpx21", "cpx11"}, Mode: "manual", Timezone: "UTC",
	})
	if err != nil {
		t.Fatalf("CreateServer: %v", err)
	}
	return srv
}

func TestAcquireActionLockSucceedsWhenFree(t *testing.T) {
	srv := newServerForAction(t)
	released, err := srv.Store().AcquireAction(srv.ID, "rescale_up", 5*time.Minute)
	if err != nil {
		t.Fatalf("AcquireAction: %v", err)
	}
	if !released {
		t.Fatal("expected released=true on first acquire")
	}
	if err := srv.Store().ReleaseAction(srv.ID); err != nil {
		t.Fatalf("ReleaseAction: %v", err)
	}
}

func TestAcquireActionLockFailsWhenBusy(t *testing.T) {
	srv := newServerForAction(t)
	if _, err := srv.Store().AcquireAction(srv.ID, "rescale_up", 5*time.Minute); err != nil {
		t.Fatalf("first AcquireAction: %v", err)
	}
	_, err := srv.Store().AcquireAction(srv.ID, "rescale_up", 5*time.Minute)
	if !errors.Is(err, ErrLocked) {
		t.Fatalf("got %v, want ErrLocked", err)
	}
}

func TestAcquireActionLockSucceedsAfterExpiry(t *testing.T) {
	srv := newServerForAction(t)
	if _, err := srv.Store().AcquireAction(srv.ID, "rescale_up", 1*time.Second); err != nil {
		t.Fatalf("first AcquireAction: %v", err)
	}
	time.Sleep(1100 * time.Millisecond)
	if _, err := srv.Store().AcquireAction(srv.ID, "rescale_up", 5*time.Minute); err != nil {
		t.Fatalf("second AcquireAction after expiry: %v", err)
	}
}