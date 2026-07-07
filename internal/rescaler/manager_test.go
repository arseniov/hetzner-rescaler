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
	var audit *store.Event
	for _, e := range events {
		if e.ID == staleID && !e.FinishedAt.IsZero() {
			pending = e
		}
		if e.Kind == "rescale_failed" && e.Error == "server restarted mid-rescale" {
			audit = e
		}
	}
	if pending == nil {
		t.Fatal("pending row's finished_at not set")
	}
	if audit == nil {
		t.Fatal("expected 1 rescale_failed recovery row, got 0")
	}
	if audit.OK {
		t.Fatalf("audit row OK = true, want false (got %+v)", audit)
	}
	if audit.TriggeredBy != "recovery" {
		t.Fatalf("audit row TriggeredBy = %q, want recovery", audit.TriggeredBy)
	}
	if !audit.StartedAt.Equal(audit.FinishedAt) {
		t.Fatalf("audit row StartedAt (%v) != FinishedAt (%v) — recovery should be instantaneous",
			audit.StartedAt, audit.FinishedAt)
	}
}

func TestStart_NoOrphansReturnsNil(t *testing.T) {
	s, _ := store.OpenTemp()
	defer s.Close()

	m := NewManager(s)
	if err := m.Start(context.Background()); err != nil {
		t.Fatalf("Start on empty store: %v", err)
	}
}

func TestStart_RecoversMultipleOrphans(t *testing.T) {
	s, _ := store.OpenTemp()
	defer s.Close()

	// Create 3 servers (each in its own project — CreateProject errors on
	// duplicate names) and seed one pending row per server.
	var srvIDs []int64
	for i := 0; i < 3; i++ {
		p, err := s.CreateProject(
			[]string{"p0", "p1", "p2"}[i],
			[]byte("tok"), []byte("nonce12byts"))
		if err != nil {
			t.Fatalf("CreateProject[%d]: %v", i, err)
		}
		srv, err := s.CreateServer(p.ID, store.Server{
			HCloudServerID: i + 1, Name: "w", Label: "w",
			BaseServerType: "cpx11", TopServerType: "cpx31",
			FallbackChain: []string{"cpx31", "cpx11"},
			Mode: "manual", Timezone: "UTC",
		})
		if err != nil {
			t.Fatalf("CreateServer[%d]: %v", i, err)
		}
		srvIDs = append(srvIDs, srv.ID)
		seedRescaleTestEvent(t, s, srv.ID, "rescale_pending")
	}

	m := NewManager(s)
	if err := m.Start(context.Background()); err != nil {
		t.Fatalf("Start: %v", err)
	}

	// Each server should now have 1 finished pending row + 1 rescale_failed audit row.
	for _, srvID := range srvIDs {
		events, _ := s.ListEventsByServer(srvID, 10)
		var pendingFinished, auditCount int
		for _, e := range events {
			if e.Kind == "rescale_pending" && !e.FinishedAt.IsZero() {
				pendingFinished++
			}
			if e.Kind == "rescale_failed" && e.Error == "server restarted mid-rescale" {
				auditCount++
			}
		}
		if pendingFinished != 1 {
			t.Errorf("server %d: pendingFinished = %d, want 1", srvID, pendingFinished)
		}
		if auditCount != 1 {
			t.Errorf("server %d: auditCount = %d, want 1", srvID, auditCount)
		}
	}
}