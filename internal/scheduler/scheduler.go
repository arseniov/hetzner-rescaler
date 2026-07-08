// Package scheduler runs the per-server goroutine that drives rescales
// based on each server's mode (scheduled, auto_promote, manual) and
// windows.
package scheduler

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/jonamat/hetzner-rescaler/internal/hetzner"
	"github.com/jonamat/hetzner-rescaler/internal/rescaler"
	"github.com/jonamat/hetzner-rescaler/internal/store"
)

// tickDebounce is the minimum interval between scheduler_tick events for the
// same server. Heartbeat events keep the operator's "last tick: Ns ago"
// staleness signal fresh without flooding the events feed.
const tickDebounce = 5 * time.Minute

// tickReasons is the closed vocabulary of scheduler_tick error fields. The
// UI maps these to human labels; tests assert exact equality.
const (
	tickReasonOKIdle         = "ok_idle"
	tickReasonNoWindows      = "no_windows"
	tickReasonAPIError       = "api_error"
	tickReasonLockContention = "lock_contention"
	tickReasonAtTarget       = "already_at_target"
)

// Scheduler holds the running goroutines for all registered servers.
type Scheduler struct {
	store     *store.Store
	api       map[int64]hetzner.API // per-project, but we use a single API for the test
	clock     Clock
	tickEvery time.Duration

	mu     sync.Mutex
	added  map[int64]struct{}
	stopCh chan struct{}
	wg     sync.WaitGroup
	log    *log.Logger
}

// New constructs a Scheduler. `tickEvery` is the interval between
// evaluations; 30s in production.
func New(st *store.Store, api hetzner.API, clk Clock, tickEvery time.Duration) *Scheduler {
	if clk == nil {
		clk = RealClock{}
	}
	return &Scheduler{
		store:     st,
		api:       map[int64]hetzner.API{0: api},
		clock:     clk,
		tickEvery: tickEvery,
		added:     map[int64]struct{}{},
		stopCh:    make(chan struct{}),
		log:       log.New(log.Writer(), "▶ scheduler ", log.LstdFlags),
	}
}

// Add registers a server to be scheduled.
func (s *Scheduler) Add(serverID int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.added[serverID]; ok {
		return
	}
	s.added[serverID] = struct{}{}
	s.wg.Add(1)
	go s.runOne(serverID)
}

// Run blocks until Stop is called.
func (s *Scheduler) Run() {
	<-s.stopCh
}

// Stop signals all goroutines to exit and waits for them.
func (s *Scheduler) Stop() {
	close(s.stopCh)
	s.wg.Wait()
}

func (s *Scheduler) runOne(serverID int64) {
	defer s.wg.Done()
	t := time.NewTicker(s.tickEvery)
	defer t.Stop()

	for {
		select {
		case <-s.stopCh:
			return
		case <-t.C:
			s.tick(serverID)
		}
	}
}

// tick does one evaluation for a server. Public for tests.
func (s *Scheduler) tick(serverID int64) {
	srv, err := s.store.GetServer(serverID)
	if err != nil {
		s.log.Printf("server %d: load failed: %v", serverID, err)
		return
	}
	api := s.apiFor(srv.ProjectID)

	switch srv.Mode {
	case "scheduled":
		s.tickScheduled(srv, api)
	case "auto_promote":
		s.tickAutoPromote(srv, api)
	case "manual":
		// no automatic rescales; housekeeping only
	}
}

func (s *Scheduler) apiFor(projectID int64) hetzner.API {
	s.mu.Lock()
	defer s.mu.Unlock()
	if a, ok := s.api[projectID]; ok {
		return a
	}
	return s.api[0]
}

func (s *Scheduler) tickScheduled(srv *store.Server, api hetzner.API) {
	wins, err := s.store.ListWindows(srv.ID)
	if err != nil {
		s.log.Printf("server %d: list windows: %v", srv.ID, err)
		return
	}
	if len(wins) == 0 {
		s.writeTickSummary(srv, tickReasonNoWindows)
		return
	}
	plain := make([]store.Window, 0, len(wins))
	for _, w := range wins {
		plain = append(plain, *w)
	}
	inWindow, target, err := EvaluateWindows(plain, s.clock.Now(), srv.Timezone)
	if err != nil {
		s.log.Printf("server %d: evaluate windows: %v", srv.ID, err)
		return
	}
	if !inWindow || target == "" {
		s.writeTickSummary(srv, tickReasonOKIdle)
		return
	}

	current, err := fetchCurrentType(context.Background(), api, srv)
	if err != nil {
		s.log.Printf("server %d: fetch current: %v", srv.ID, err)
		s.writeTickSummary(srv, tickReasonAPIError)
		return
	}
	if current == target {
		s.writeTickSummary(srv, tickReasonAtTarget)
		return
	}
	s.dispatch(srv, api, target, "scheduler")
}

func (s *Scheduler) tickAutoPromote(srv *store.Server, api hetzner.API) {
	if srv.PromoteState == nil {
		return
	}
	current, err := fetchCurrentType(context.Background(), api, srv)
	if err != nil {
		s.log.Printf("server %d: fetch current: %v", srv.ID, err)
		s.writeTickSummary(srv, tickReasonAPIError)
		return
	}
	switch *srv.PromoteState {
	case "promote_requested":
		if current == srv.BaseServerType {
			s.dispatch(srv, api, srv.TopServerType, "auto_promote")
			return
		}
		// Gate unmet — surface it. (current == top means we're already where
		// the operator asked; current elsewhere means we're between states.)
		s.writeTickSummary(srv, tickReasonOKIdle)
	case "demote_requested":
		if current == srv.TopServerType {
			s.dispatch(srv, api, srv.BaseServerType, "manual")
			return
		}
		s.writeTickSummary(srv, tickReasonOKIdle)
	}
}

// writeTickSummary persists a scheduler_tick event for srv unless one was
// written within the last tickDebounce. The event's Error field carries a
// small closed-vocabulary reason — see the tickReason* constants.
//
// No-op if reason is empty (the caller has nothing to surface to the
// operator).
func (s *Scheduler) writeTickSummary(srv *store.Server, reason string) {
	if reason == "" {
		return
	}
	now := s.clock.Now().UTC()
	last, err := s.store.LastEventOfKind(srv.ID, "scheduler_tick")
	if err != nil {
		s.log.Printf("server %d: read last tick: %v", srv.ID, err)
		return
	}
	if last.ID != 0 && now.Sub(last.StartedAt) < tickDebounce {
		return
	}
	if _, err := s.store.AppendEvent(store.Event{
		ServerID:    srv.ID,
		Kind:        "scheduler_tick",
		StartedAt:   now,
		FinishedAt:  now,
		OK:          true,
		TriggeredBy: "scheduler",
		Error:       reason,
	}); err != nil {
		s.log.Printf("server %d: write scheduler_tick (%s): %v", srv.ID, reason, err)
	}
}

func (s *Scheduler) tryPromote(srv *store.Server, api hetzner.API, target string) error {
	// Try the rescale; on unavailable, set promote_state to 'promoting' so we
	// keep trying on future ticks. The actual polling cadence is the
	// scheduler's tick interval; no separate goroutine.
	hsrv, err := api.GetServer(context.Background(), srv.HCloudServerID)
	if err != nil {
		return err
	}
	used, err := rescaler.RescaleWithFallback(context.Background(), api, hsrv, target, srv.FallbackChain)
	_ = used
	return err
}

// dispatch locks the action, runs the rescale, writes an event, releases.
func (s *Scheduler) dispatch(srv *store.Server, api hetzner.API, target, triggeredBy string) {
	acquired, err := s.store.AcquireAction(srv.ID, "rescale_to_"+target, 30*time.Minute)
	if err != nil {
		s.log.Printf("server %d: acquire action: %v", srv.ID, err)
		if errors.Is(err, store.ErrLocked) {
			s.writeTickSummary(srv, tickReasonLockContention)
		}
		return
	}
	if !acquired {
		s.log.Printf("server %d: another action in flight, skipping", srv.ID)
		s.writeTickSummary(srv, tickReasonLockContention)
		return
	}
	defer s.store.ReleaseAction(srv.ID)

	hsrv, err := api.GetServer(context.Background(), srv.HCloudServerID)
	if err != nil {
		s.log.Printf("server %d: fetch server: %v", srv.ID, err)
		return
	}
	current := ""
	if hsrv.ServerType != nil {
		current = hsrv.ServerType.Name
	}

	start := time.Now().UTC()
	used, rescaleErr := rescaler.RescaleWithFallback(context.Background(), api, hsrv, target, srv.FallbackChain)
	finished := time.Now().UTC()

	kind := "rescale_up"
	if used == srv.BaseServerType {
		kind = "rescale_down"
	}
	if rescaleErr != nil {
		if _, err := s.store.AppendEvent(store.Event{
			ServerID:    srv.ID,
			Kind:        "rescale_failed",
			FromType:    current,
			ToType:      target,
			StartedAt:   start,
			FinishedAt:  finished,
			OK:          false,
			Error:       rescaleErr.Error(),
			TriggeredBy: triggeredBy,
		}); err != nil {
			s.log.Printf("server %d: write failed event: %v", srv.ID, err)
		}
		return
	}
	if _, err := s.store.AppendEvent(store.Event{
		ServerID:    srv.ID,
		Kind:        kind,
		FromType:    current,
		ToType:      used,
		StartedAt:   start,
		FinishedAt:  finished,
		OK:          true,
		TriggeredBy: triggeredBy,
	}); err != nil {
		s.log.Printf("server %d: write event: %v", srv.ID, err)
	}
}

func fetchCurrentType(ctx context.Context, api hetzner.API, srv *store.Server) (string, error) {
	s, err := api.GetServer(ctx, srv.HCloudServerID)
	if err != nil {
		return "", err
	}
	if s == nil || s.ServerType == nil {
		return "", nil
	}
	return s.ServerType.Name, nil
}