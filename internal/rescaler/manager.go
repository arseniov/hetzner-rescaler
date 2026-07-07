package rescaler

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/jonamat/hetzner-rescaler/internal/hetzner"
	"github.com/jonamat/hetzner-rescaler/internal/store"
)

// ErrAlreadyInProgress is returned by Manager.Submit when a goroutine for
// the given server is already running. The handler maps this to a 409.
var ErrAlreadyInProgress = errors.New("rescaler: rescale already in progress")

// Manager owns the lifecycle of async rescale goroutines.
type Manager struct {
	store       *store.Store
	mu          sync.Mutex
	jobs        map[int64]context.CancelFunc
	apiResolver apiResolver
}

func NewManager(s *store.Store) *Manager {
	return &Manager{store: s, jobs: make(map[int64]context.CancelFunc)}
}

// Start scans for orphaned rescale_pending rows from a previous process
// run and marks each as failed. Idempotent — safe to call on every boot.
// Must be called once before any Submit.
//
// Partial-failure semantics: if AppendEvent or UpdateEventFinished fails
// mid-loop, Start returns the error and any remaining orphans are left for
// the next boot to retry. The already-written audit rows remain, which may
// produce duplicate rescale_failed rows on retry. This is acceptable for a
// boot-time recovery operation: the operator sees at most one extra audit
// row, the pending row does eventually close out, and no data is lost.
func (m *Manager) Start(ctx context.Context) error {
	rows, err := m.store.DB().QueryContext(ctx,
		`SELECT id, server_id FROM events
		 WHERE kind = 'rescale_pending' AND finished_at IS NULL`,
	)
	if err != nil {
		return err
	}
	defer rows.Close()
	type orphan struct {
		id       int64
		serverID int64
	}
	var orphans []orphan
	for rows.Next() {
		var o orphan
		if err := rows.Scan(&o.id, &o.serverID); err != nil {
			return err
		}
		orphans = append(orphans, o)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	now := time.Now().UTC()
	for _, o := range orphans {
		if _, err := m.store.AppendEvent(store.Event{
			ServerID:    o.serverID,
			Kind:        "rescale_failed",
			StartedAt:   now,
			FinishedAt:  now,
			OK:          false,
			Error:       "server restarted mid-rescale",
			TriggeredBy: "recovery",
		}); err != nil {
			return err
		}
		if err := m.store.UpdateEventFinished(o.id, false, "server restarted mid-rescale"); err != nil {
			return err
		}
	}
	return nil
}

// Submit reserves the server for an async rescale. Returns the new
// pending event ID, or ErrAlreadyInProgress if the server already has a
// goroutine in flight. The actual rescale work is not started here —
// Task 10 introduces the goroutine + runRescale. This task only proves
// the row insertion + concurrency check.
func (m *Manager) Submit(ctx context.Context, srv *store.Server, target string, triggeredBy string) (int64, error) {
	m.mu.Lock()
	if _, busy := m.jobs[srv.ID]; busy {
		m.mu.Unlock()
		return 0, ErrAlreadyInProgress
	}
	m.mu.Unlock()

	api, err := m.resolveAPI(ctx, srv.ProjectID)
	if err != nil {
		return 0, err
	}
	hserver, err := api.GetServer(ctx, srv.HCloudServerID)
	if err != nil {
		return 0, err
	}
	fromType := ""
	if hserver != nil && hserver.ServerType != nil {
		fromType = hserver.ServerType.Name
	}
	now := time.Now().UTC()
	id, err := m.store.AppendEvent(store.Event{
		ServerID:    srv.ID,
		Kind:        "rescale_pending",
		FromType:    fromType,
		StartedAt:   now,
		TriggeredBy: triggeredBy,
	})
	if err != nil {
		return 0, err
	}

	// Register the job in the map. Use a no-op cancel for now; Task 10
	// replaces it with the goroutine's real cancel function.
	m.mu.Lock()
	m.jobs[srv.ID] = func() {}
	m.mu.Unlock()

	return id, nil
}

// resolveAPI returns a hetzner.API for the given project. Wired in
// Task 14 when the Manager is integrated into Deps; for now it's a
// hook so tests can stub via setAPIResolver.
func (m *Manager) resolveAPI(ctx context.Context, projectID int64) (hetzner.API, error) {
	if m.apiResolver == nil {
		return nil, errors.New("rescaler: api resolver not configured")
	}
	return m.apiResolver(ctx, projectID)
}

// apiResolver is set by the cmd wiring. Tests override via setAPIResolver.
type apiResolver func(ctx context.Context, projectID int64) (hetzner.API, error)

func (m *Manager) setAPIResolver(fn apiResolver) { m.apiResolver = fn }
