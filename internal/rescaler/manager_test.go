package rescaler

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/jonamat/hetzner-rescaler/internal/hcloudmock"
	"github.com/jonamat/hetzner-rescaler/internal/hetzner"
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

func TestSubmit_InsertsPendingRowAndReturnsID(t *testing.T) {
	s, _ := store.OpenTemp()
	defer s.Close()
	_, srvID := seedRescaleTestProjectAndServer(t, s)
	srv, _ := s.GetServer(srvID)

	m := NewManager(s)
	if err := m.Start(context.Background()); err != nil {
		t.Fatalf("Start: %v", err)
	}
	mock := hcloudmock.New()
	mock.AddServer(&hetzner.Server{ID: int(srv.HCloudServerID), ServerType: &hetzner.ServerType{Name: "cpx11"}})
	m.setAPIResolver(func(ctx context.Context, projectID int64) (hetzner.API, error) {
		return mock, nil
	})
	id, err := m.Submit(context.Background(), srv, "cpx31", "api")
	if err != nil {
		t.Fatalf("Submit: %v", err)
	}
	if id == 0 {
		t.Fatal("expected non-zero pending event ID")
	}
	events, _ := s.ListEventsByServer(srvID, 10)
	var pending *store.Event
	for _, e := range events {
		if e.ID == id && e.Kind == "rescale_pending" {
			pending = e
		}
	}
	if pending == nil {
		t.Fatalf("pending row not written: %+v", events)
	}
	if pending.FromType != "cpx11" {
		t.Fatalf("FromType = %q, want cpx11 (current)", pending.FromType)
	}
	if !pending.FinishedAt.IsZero() {
		t.Fatalf("FinishedAt should be zero, got %v", pending.FinishedAt)
	}
}

func TestSubmit_RejectsWhenAlreadyInProgress(t *testing.T) {
	s, _ := store.OpenTemp()
	defer s.Close()
	_, srvID := seedRescaleTestProjectAndServer(t, s)
	srv, _ := s.GetServer(srvID)

	m := NewManager(s)
	_ = m.Start(context.Background())
	m.setAPIResolver(func(ctx context.Context, projectID int64) (hetzner.API, error) {
		return hcloudmock.New(), nil
	})

	// Manually mark the server busy (the goroutine hasn't been introduced yet).
	m.mu.Lock()
	m.jobs[srvID] = func() {}
	m.mu.Unlock()

	_, err := m.Submit(context.Background(), srv, "cpx31", "api")
	if !errors.Is(err, ErrAlreadyInProgress) {
		t.Fatalf("err = %v, want ErrAlreadyInProgress", err)
	}
}

func TestRunRescale_HappyPathWritesTerminalWithReconciledToType(t *testing.T) {
	s, _ := store.OpenTemp()
	defer s.Close()
	_, srvID := seedRescaleTestProjectAndServer(t, s)
	srv, _ := s.GetServer(srvID)

	mock := hcloudmock.New()
	mock.AddServer(&hetzner.Server{
		ID: srv.HCloudServerID, Name: srv.Name,
		ServerType: &hetzner.ServerType{Name: "cpx11"},
	})

	m := NewManager(s)
	_ = m.Start(context.Background())
	m.setAPIResolver(func(ctx context.Context, projectID int64) (hetzner.API, error) {
		return mock, nil
	})

	done := make(chan struct{})
	go func() {
		_, _ = m.Submit(context.Background(), srv, "cpx31", "api")
		close(done)
	}()
	<-done

	// Wait for the goroutine to finish (poll the jobs map).
	// Mock actions take ~5s each (pollInterval) — so the rescale takes
	// at least ~10s end-to-end (change-type + power-on). Poll for
	// 60s to leave headroom for slow CI.
	deadline := time.Now().Add(60 * time.Second)
	for time.Now().Before(deadline) {
		m.mu.Lock()
		_, busy := m.jobs[srvID]
		m.mu.Unlock()
		if !busy {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}

	events, _ := s.ListEventsByServer(srvID, 20)
	var pending *store.Event
	var terminal *store.Event
	for _, e := range events {
		if e.Kind == "rescale_pending" {
			pending = e
		}
		if e.Kind == "rescale_completed" || e.Kind == "rescale_failed" {
			terminal = e
		}
	}
	if pending == nil {
		t.Fatal("pending row missing")
	}
	if pending.FinishedAt.IsZero() {
		t.Fatalf("pending row should be finished: %+v", pending)
	}
	if terminal == nil {
		t.Fatalf("terminal row missing; events=%+v", events)
	}
	if terminal.ToType != "cpx31" {
		t.Fatalf("terminal ToType = %q, want cpx31 (reconciled)", terminal.ToType)
	}
	if !terminal.OK {
		t.Fatalf("terminal row ok = false, want true: %+v", terminal)
	}
}

func TestRunRescale_FailureWritesFailedTerminal(t *testing.T) {
	s, _ := store.OpenTemp()
	defer s.Close()
	_, srvID := seedRescaleTestProjectAndServer(t, s)
	srv, _ := s.GetServer(srvID)

	mock := hcloudmock.New()
	mock.AddServer(&hetzner.Server{
		ID: srv.HCloudServerID, Name: srv.Name,
		ServerType: &hetzner.ServerType{Name: "cpx11"},
	})
	mock.SetChangeTypeOverride(func(target *hetzner.ServerType) error {
		return errors.New("simulated change_type failure")
	})

	m := NewManager(s)
	_ = m.Start(context.Background())
	m.setAPIResolver(func(ctx context.Context, projectID int64) (hetzner.API, error) {
		return mock, nil
	})

	_, _ = m.Submit(context.Background(), srv, "cpx31", "api")

	// Wait for the goroutine to finish.
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		m.mu.Lock()
		_, busy := m.jobs[srvID]
		m.mu.Unlock()
		if !busy {
			break
		}
		time.Sleep(20 * time.Millisecond)
	}

	events, _ := s.ListEventsByServer(srvID, 20)
	var terminal *store.Event
	for _, e := range events {
		if e.Kind == "rescale_failed" {
			terminal = e
		}
	}
	if terminal == nil {
		t.Fatalf("expected rescale_failed row, got %+v", events)
	}
	if terminal.OK {
		t.Fatal("terminal OK should be false")
	}
	if !strings.Contains(terminal.Error, "simulated change_type failure") {
		t.Fatalf("terminal error = %q, want it to mention simulated failure", terminal.Error)
	}
}