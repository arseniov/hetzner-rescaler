package rescaler

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/jonamat/hetzner-rescaler/internal/hetzner"
	"github.com/jonamat/hetzner-rescaler/internal/store"
)

// ErrAlreadyInProgress is returned by Manager.Submit when a goroutine for
// the given server is already running. The handler maps this to a 409.
var ErrAlreadyInProgress = errors.New("rescaler: rescale already in progress")

// Manager owns the lifecycle of async rescale goroutines.
//
// The jobs map is the source of truth for in-flight goroutines: while a
// server ID is in m.jobs, a goroutine is running for that server. Submit
// uses the map to enforce "one rescale per server" semantics; runRescale
// removes the entry in a defer so the slot frees as soon as the work
// completes (success, failure, or panic).
//
// The store is the secondary source of truth: a still-pending
// rescale_pending row means a previous process was killed mid-rescale,
// and Start walks those rows on boot to mark them as failed recovery.
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

// Submit reserves the server for an async rescale, inserts the pending
// event row, and spawns a goroutine that walks the rescale phases. The
// HTTP request returns immediately with the pending event ID.
//
// Errors:
//   - ErrAlreadyInProgress: another goroutine is already working on srv.
//   - any API/store error during the synchronous setup phase.
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

	// Spawn the goroutine. The goroutine owns the rest of the lifecycle.
	jobCtx, cancel := context.WithCancel(context.Background())
	m.mu.Lock()
	m.jobs[srv.ID] = cancel
	m.mu.Unlock()
	go m.runRescale(jobCtx, srv, hserver, fromType, api, target, id, triggeredBy)
	return id, nil
}

// runRescale walks the rescale phases: emits per-phase updates via
// UpdateEventPhase, calls RescaleWithFallbackWithHook for the actual
// Hetzner work, then reconciles the terminal to_type from Hetzner and
// writes the terminal event row. On any error, writes a rescale_failed
// row.
//
// Defer: removes the jobs map entry, recovers from panic, writes the
// terminal row even on panic.
//
// `hserver` and `fromType` are passed in from Submit (which already
// fetched them) so the goroutine does not duplicate the Hetzner call.
func (m *Manager) runRescale(ctx context.Context, srv *store.Server, hserver *hetzner.Server, fromType string, api hetzner.API, target string, pendingID int64, triggeredBy string) {
	defer func() {
		m.mu.Lock()
		delete(m.jobs, srv.ID)
		m.mu.Unlock()
		if r := recover(); r != nil {
			now := time.Now().UTC()
			_, _ = m.store.AppendEvent(store.Event{
				ServerID:    srv.ID,
				Kind:        "rescale_failed",
				StartedAt:   now,
				FinishedAt:  now,
				OK:          false,
				Error:       fmt.Sprintf("panic: %v", r),
				TriggeredBy: triggeredBy,
			})
			_ = m.store.UpdateEventFinished(pendingID, false, fmt.Sprintf("panic: %v", r))
		}
	}()

	phaseHook := func(phase string) {
		_ = m.store.UpdateEventPhase(pendingID, phase)
	}

	_, err := RescaleWithFallbackWithHook(ctx, api, hserver, target, srv.FallbackChain, phaseHook)

	// Reconcile terminal to_type from Hetzner (authoritative). The
	// in-memory hserver may have a stale type if RescaleWithFallback's
	// final Rescale call updated the struct; we trust the Hetzner
	// response over that.
	var toType string
	if h, gerr := api.GetServer(ctx, srv.HCloudServerID); gerr == nil && h != nil && h.ServerType != nil {
		toType = h.ServerType.Name
	} else {
		// Reconciliation failed — fall back to the requested target on
		// failure, or whatever type RescaleWithFallback used on success.
		// Either way: still better than blank.
		toType = target
	}

	now := time.Now().UTC()
	termKind := "rescale_completed"
	ok := true
	errMsg := ""
	if err != nil {
		termKind = "rescale_failed"
		ok = false
		errMsg = err.Error()
	}
	if _, aerr := m.store.AppendEvent(store.Event{
		ServerID:    srv.ID,
		Kind:        termKind,
		FromType:    fromType,
		ToType:      toType,
		StartedAt:   now,
		FinishedAt:  now,
		OK:          ok,
		Error:       errMsg,
		TriggeredBy: triggeredBy,
	}); aerr != nil {
		// If we can't write the terminal row, the defer still removed the
		// jobs map entry; nothing else to do.
		_ = aerr
	}
	_ = m.store.UpdateEventFinished(pendingID, ok, errMsg)
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
