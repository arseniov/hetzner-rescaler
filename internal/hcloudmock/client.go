// Package hcloudmock is an in-memory fake of the hetzner.API for tests.
package hcloudmock

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/jonamat/hetzner-rescaler/internal/hetzner"
)

// Fake implements hetzner.API with in-memory state. Safe for concurrent use
// across the goroutines a test might spawn.
type Fake struct {
	mu sync.Mutex

	servers       map[int64]*hetzner.Server
	serverTypes   map[string]*hetzner.ServerType
	actions       map[int64]*hetzner.Action
	nextActionID  int64
	unavailable   map[string]bool
	shutdownCount int

	// errorOverrides: if a function is set, it returns its result for the
	// corresponding (server, type) pair. Used by tests that need a custom
	// failure shape.
	changeTypeOverride func(target *hetzner.ServerType) error
}

// New returns a fresh Fake with a few default server types pre-populated.
func New() *Fake {
	f := &Fake{
		servers:     map[int64]*hetzner.Server{},
		serverTypes: map[string]*hetzner.ServerType{},
		actions:     map[int64]*hetzner.Action{},
		unavailable: map[string]bool{},
	}
	for _, n := range []string{"cpx11", "cpx21", "cpx31", "cx11", "cx21", "cx31"} {
		f.serverTypes[n] = &hetzner.ServerType{Name: n}
	}
	return f
}

// AddServer registers a server in the fake. If s.ID is 0, AddServer assigns
// a new ID and writes it back into the caller's *Server — this mutates
// the input.
func (f *Fake) AddServer(s *hetzner.Server) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if s.ID == 0 {
		s.ID = int64(len(f.servers) + 1)
	}
	f.servers[s.ID] = s
}

// MarkUnavailable marks a server type as out of stock. Future
// ChangeServerType calls targeting that type return ErrUnavailable.
func (f *Fake) MarkUnavailable(name string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.unavailable[name] = true
}

// SetChangeTypeOverride lets a test inject a custom error from
// ChangeServerType. Used to verify the error-classification path.
func (f *Fake) SetChangeTypeOverride(fn func(target *hetzner.ServerType) error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.changeTypeOverride = fn
}

// ShutdownCount returns the number of ShutdownServer calls the fake has
// received. Tests use it to verify that a flow didn't induce an extra
// shutdown — e.g. when the pre-check aborts early because every
// candidate type is out of stock, a running server should never have
// been asked to shut down at all.
func (f *Fake) ShutdownCount() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.shutdownCount
}

// ---- hetzner.API ----

func (f *Fake) GetServer(_ context.Context, id int64) (*hetzner.Server, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	s, ok := f.servers[id]
	if !ok {
		return nil, fmt.Errorf("hcloudmock: server %d not found", id)
	}
	return s, nil
}

func (f *Fake) ListServers(_ context.Context) ([]*hetzner.Server, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	out := make([]*hetzner.Server, 0, len(f.servers))
	for _, s := range f.servers {
		out = append(out, s)
	}
	return out, nil
}

func (f *Fake) ShutdownServer(_ context.Context, srv *hetzner.Server) (*hetzner.Action, error) {
	f.mu.Lock()
	f.shutdownCount++
	f.mu.Unlock()
	return f.startAction(srv, "shutdown")
}

func (f *Fake) ChangeServerType(_ context.Context, srv *hetzner.Server, target *hetzner.ServerType) (*hetzner.Action, error) {
	f.mu.Lock()
	if f.changeTypeOverride != nil {
		err := f.changeTypeOverride(target)
		f.mu.Unlock()
		if err != nil {
			return nil, err
		}
		// fall through to success
	} else if f.unavailable[target.Name] {
		f.mu.Unlock()
		return nil, fmt.Errorf("hcloudmock: %w: %s", hetzner.ErrUnavailable, target.Name)
	} else {
		f.mu.Unlock()
	}
	return f.startAction(srv, "change_type")
}

func (f *Fake) PowerOnServer(_ context.Context, srv *hetzner.Server) (*hetzner.Action, error) {
	return f.startAction(srv, "poweron")
}

func (f *Fake) ListServerTypes(_ context.Context) ([]*hetzner.ServerType, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	out := make([]*hetzner.ServerType, 0, len(f.serverTypes))
	for _, t := range f.serverTypes {
		out = append(out, t)
	}
	return out, nil
}

func (f *Fake) GetServerType(_ context.Context, name string) (*hetzner.ServerType, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	t, ok := f.serverTypes[name]
	if !ok {
		return nil, fmt.Errorf("hcloudmock: type %q not found", name)
	}
	// Return a copy so callers can't mutate the fake's internal
	// state. Hetzner's hcloud.ServerType doesn't expose availability
	// (the API only reports it implicitly via ChangeServerType
	// failures), so there's nothing to surface here — the test
	// harness asserts "unavailable" by failing ChangeServerType.
	out := *t
	return &out, nil
}

func (f *Fake) GetAction(_ context.Context, id int64) (*hetzner.Action, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	a, ok := f.actions[id]
	if !ok {
		return nil, fmt.Errorf("hcloudmock: action %d not found", id)
	}
	// Simulate progress: every call advances the action. The first two
	// calls return "running", the third returns "success". This mirrors
	// Hetzner's typical behavior for short actions.
	switch a.Progress {
	case 0:
		a.Progress = 50
		a.Status = hetzner.ActionStatusRunning
	case 50:
		a.Progress = 100
		a.Status = hetzner.ActionStatusSuccess
		a.Finished = time.Now().UTC()
	}
	return a, nil
}

// ---- helpers ----

func (f *Fake) startAction(srv *hetzner.Server, command string) (*hetzner.Action, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.nextActionID++
	a := &hetzner.Action{
		ID:       f.nextActionID,
		Status:   hetzner.ActionStatusRunning,
		Progress: 0,
		Started:  time.Now().UTC(),
		Command:  command,
	}
	a.Resources = []*hcloud.ActionResource{{Type: "server", ID: srv.ID}}
	f.actions[a.ID] = a
	return a, nil
}
