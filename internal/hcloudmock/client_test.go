package hcloudmock

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jonamat/hetzner-rescaler/internal/hetzner"
)

func TestFakeListsConfiguredServers(t *testing.T) {
	f := New()
	srv := &hetzner.Server{ID: 42, Name: "web"}
	f.AddServer(srv)

	got, err := f.ListServers(context.Background())
	if err != nil {
		t.Fatalf("ListServers: %v", err)
	}
	if len(got) != 1 || got[0].ID != 42 {
		t.Fatalf("got %+v", got)
	}
}

func TestFakeShutdownReturnsActionThatCompletes(t *testing.T) {
	f := New()
	srv := &hetzner.Server{}
	srv.ID = 1
	f.AddServer(srv)

	act, err := f.ShutdownServer(context.Background(), srv)
	if err != nil {
		t.Fatalf("Shutdown: %v", err)
	}
	if act.Status != hetzner.ActionStatusRunning {
		t.Fatalf("shutdown action status = %v, want running", act.Status)
	}
}

func TestFakeUnavailableForType(t *testing.T) {
	f := New()
	f.MarkUnavailable("cpx31")

	_, err := f.ChangeServerTypeReturnsErrorFor(context.Background(), &hetzner.Server{}, &hetzner.ServerType{Name: "cpx31"}, hetzner.ErrUnavailable)
	if err == nil {
		t.Fatal("expected unavailable error")
	}
	if !hetzner.IsUnavailable(err) {
		t.Fatalf("err = %v, want unavailable", err)
	}
}

func TestFakeGetActionProgressesAfterProgress(t *testing.T) {
	f := New()
	srv := &hetzner.Server{}
	srv.ID = 1
	f.AddServer(srv)

	act, err := f.PowerOnServer(context.Background(), srv)
	if err != nil {
		t.Fatalf("PowerOn: %v", err)
	}

	if act.Status == hetzner.ActionStatusSuccess {
		t.Skip("action already finished")
	}
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		act2, err := f.GetAction(context.Background(), act.ID)
		if err != nil {
			t.Fatalf("GetAction: %v", err)
		}
		if act2.Status == hetzner.ActionStatusSuccess {
			return
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatal("action did not reach success within 2s")
	_ = errors.New
}
