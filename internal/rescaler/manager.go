package rescaler

import (
	"context"
	"sync"
	"time"

	"github.com/jonamat/hetzner-rescaler/internal/store"
)

// ErrAlreadyInProgress is returned by Manager.Submit when a goroutine for
// the given server is already running. The handler maps this to a 409
// response with the pending event ID so the UI can link to it.
var ErrAlreadyInProgress = context.DeadlineExceeded // placeholder; replaced in Task 8

// Manager owns the lifecycle of async rescale goroutines. One instance
// is shared across the API and the scheduler (the scheduler is the same
// process in `serve` mode; the API is the only consumer).
//
// Concurrency model:
//   - jobs is the single source of truth for "is this server mid-rescale?"
//   - Each entry holds the goroutine's cancel function so Shutdown can
//     cancel all in-flight jobs gracefully.
//   - The store is the secondary source of truth (via the rescale_pending
//     row) so a process restart can recover orphaned jobs.
type Manager struct {
	store *store.Store

	mu   sync.Mutex
	jobs map[int64]context.CancelFunc // serverID -> cancel
}

// NewManager constructs a Manager bound to the given store.
func NewManager(s *store.Store) *Manager {
	return &Manager{
		store: s,
		jobs:  make(map[int64]context.CancelFunc),
	}
}

// Start scans for orphaned rescale_pending rows from a previous process
// run and marks each as failed. Idempotent — safe to call on every boot.
// Must be called once before any Submit.
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
		// Append a rescale_failed audit row first so the operator sees a
		// terminal event with a clear cause.
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
		// Then close out the pending row.
		if err := m.store.UpdateEventFinished(o.id, false, "server restarted mid-rescale"); err != nil {
			return err
		}
	}
	return nil
}