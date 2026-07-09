package store

import (
	"testing"
	"time"

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

// TestCreateServerBroadcastsCreated verifies that CreateServer emits a
// "created" lifecycle event with the new server's ID, when a hub is attached.
func TestCreateServerBroadcastsCreated(t *testing.T) {
	s := newTestStore(t)
	p, _ := s.CreateProject("p", []byte("tok"), []byte("nonce12byts"))

	hub := broadcast.NewHub[ServerLifecycleEvent]()
	ch, unsub := hub.Subscribe(4)
	defer unsub()
	s.SetServerLifecycleHub(hub)

	srv, err := s.CreateServer(p.ID, Server{
		HCloudServerID: 1, Name: "web", BaseServerType: "cpx11", TopServerType: "cpx21",
		FallbackChain: []string{"cpx21"}, Mode: "manual", Timezone: "UTC",
	})
	if err != nil {
		t.Fatalf("CreateServer: %v", err)
	}

	select {
	case got := <-ch:
		if got.Kind != "created" {
			t.Fatalf("Kind = %q, want %q", got.Kind, "created")
		}
		if got.ServerID != srv.ID {
			t.Fatalf("ServerID = %d, want %d", got.ServerID, srv.ID)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for created event")
	}
}

// TestUpdateServerBroadcastsUpdated verifies UpdateServer emits an
// "updated" event with the same server ID.
func TestUpdateServerBroadcastsUpdated(t *testing.T) {
	s := newTestStore(t)
	p, _ := s.CreateProject("p", []byte("tok"), []byte("nonce12byts"))
	srv, _ := s.CreateServer(p.ID, Server{
		HCloudServerID: 1, Name: "web", BaseServerType: "cpx11", TopServerType: "cpx21",
		FallbackChain: []string{"cpx21"}, Mode: "manual", Timezone: "UTC",
	})

	hub := broadcast.NewHub[ServerLifecycleEvent]()
	ch, unsub := hub.Subscribe(4)
	defer unsub()
	s.SetServerLifecycleHub(hub)

	srv.Label = "updated"
	if err := s.UpdateServer(*srv); err != nil {
		t.Fatalf("UpdateServer: %v", err)
	}

	select {
	case got := <-ch:
		if got.Kind != "updated" || got.ServerID != srv.ID {
			t.Fatalf("got %+v, want {Kind: updated, ServerID: %d}", got, srv.ID)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for updated event")
	}
}

// TestDeleteServerBroadcastsDeleted verifies DeleteServer emits a "deleted"
// event with the deleted server ID.
func TestDeleteServerBroadcastsDeleted(t *testing.T) {
	s := newTestStore(t)
	p, _ := s.CreateProject("p", []byte("tok"), []byte("nonce12byts"))
	srv, _ := s.CreateServer(p.ID, Server{
		HCloudServerID: 1, Name: "web", BaseServerType: "cpx11", TopServerType: "cpx21",
		FallbackChain: []string{"cpx21"}, Mode: "manual", Timezone: "UTC",
	})

	hub := broadcast.NewHub[ServerLifecycleEvent]()
	ch, unsub := hub.Subscribe(4)
	defer unsub()
	s.SetServerLifecycleHub(hub)

	if err := s.DeleteServer(srv.ID); err != nil {
		t.Fatalf("DeleteServer: %v", err)
	}

	select {
	case got := <-ch:
		if got.Kind != "deleted" || got.ServerID != srv.ID {
			t.Fatalf("got %+v, want {Kind: deleted, ServerID: %d}", got, srv.ID)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for deleted event")
	}
}

// TestMutatorsNoHubIsNoop verifies the mutators succeed without a hub.
func TestMutatorsNoHubIsNoop(t *testing.T) {
	s := newTestStore(t)
	p, _ := s.CreateProject("p", []byte("tok"), []byte("nonce12byts"))
	srv, _ := s.CreateServer(p.ID, Server{
		HCloudServerID: 1, Name: "web", BaseServerType: "cpx11", TopServerType: "cpx21",
		FallbackChain: []string{"cpx21"}, Mode: "manual", Timezone: "UTC",
	})
	// No SetServerLifecycleHub call; s.serverLifecycleHub is nil.
	srv.Label = "x"
	if err := s.UpdateServer(*srv); err != nil {
		t.Fatalf("UpdateServer: %v", err)
	}
	if err := s.DeleteServer(srv.ID); err != nil {
		t.Fatalf("DeleteServer: %v", err)
	}
}
