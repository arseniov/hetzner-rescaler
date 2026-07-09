package store

import (
	"testing"

	"github.com/jonamat/hetzner-rescaler/internal/broadcast"
)

// TestSetServerLifecycleHub_AccessorRoundTrip verifies the setter stores
// the hub and the getter returns the same pointer.
func TestSetServerLifecycleHub_AccessorRoundTrip(t *testing.T) {
	s := newTestStore(t)
	hub := broadcast.NewHub[ServerLifecycleEvent]()
	s.SetServerLifecycleHub(hub)
	if got := s.ServerLifecycleHub(); got != hub {
		t.Fatalf("ServerLifecycleHub() = %p, want %p", got, hub)
	}
}

// TestSetServerLifecycleHub_DetachWithNil verifies passing nil detaches
// the hub — subsequent getter returns nil.
func TestSetServerLifecycleHub_DetachWithNil(t *testing.T) {
	s := newTestStore(t)
	hub := broadcast.NewHub[ServerLifecycleEvent]()
	s.SetServerLifecycleHub(hub)
	s.SetServerLifecycleHub(nil)
	if got := s.ServerLifecycleHub(); got != nil {
		t.Fatalf("ServerLifecycleHub() after detach = %p, want nil", got)
	}
}