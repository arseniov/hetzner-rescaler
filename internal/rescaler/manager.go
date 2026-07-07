package rescaler

import (
	"context"
	"sync"

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