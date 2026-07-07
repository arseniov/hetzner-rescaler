package rescaler

import (
	"context"
	"testing"
	"time"

	"github.com/jonamat/hetzner-rescaler/internal/store"
)

func TestNewManager_NotNil(t *testing.T) {
	s, err := store.OpenTemp()
	if err != nil {
		t.Fatalf("OpenTemp: %v", err)
	}
	defer s.Close()
	m := NewManager(s)
	if m == nil {
		t.Fatal("NewManager returned nil")
	}
	if m.jobs == nil {
		t.Fatal("jobs map not initialised")
	}
}

func seedRescaleTestProjectAndServer(t *testing.T, s *store.Store) (int64, int64) {
	t.Helper()
	p, err := s.CreateProject("p", []byte("tok"), []byte("nonce12byts"))
	if err != nil {
		t.Fatalf("CreateProject: %v", err)
	}
	srv, err := s.CreateServer(p.ID, store.Server{
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

func seedRescaleTestEvent(t *testing.T, s *store.Store, serverID int64, kind string) int64 {
	t.Helper()
	id, err := s.AppendEvent(store.Event{
		ServerID: serverID, Kind: kind,
		StartedAt: time.Now().UTC(), TriggeredBy: "test",
	})
	if err != nil {
		t.Fatalf("AppendEvent: %v", err)
	}
	return id
}

func TestStart_RecoversOrphanPendingRows(t *testing.T) {
	s, err := store.OpenTemp()
	if err != nil {
		t.Fatalf("OpenTemp: %v", err)
	}
	defer s.Close()
	_, srvID := seedRescaleTestProjectAndServer(t, s)

	// Insert a stale pending row from a "previous run".
	staleID := seedRescaleTestEvent(t, s, srvID, "rescale_pending")

	m := NewManager(s)
	if err := m.Start(context.Background()); err != nil {
		t.Fatalf("Start: %v", err)
	}

	// Original pending row should now have finished_at set.
	events, _ := s.ListEventsByServer(srvID, 10)
	var pending *store.Event
	var failedCount int
	for _, e := range events {
		if e.ID == staleID && !e.FinishedAt.IsZero() {
			pending = e
		}
		if e.Kind == "rescale_failed" && e.Error == "server restarted mid-rescale" {
			failedCount++
		}
	}
	if pending == nil {
		t.Fatal("pending row's finished_at not set")
	}
	if failedCount != 1 {
		t.Fatalf("expected 1 rescale_failed recovery row, got %d", failedCount)
	}
}