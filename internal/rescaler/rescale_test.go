package rescaler

import (
	"context"
	"testing"
	"time"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/jonamat/hetzner-rescaler/internal/hcloudmock"
	"github.com/jonamat/hetzner-rescaler/internal/hetzner"
)

func TestRescaleToTargetType(t *testing.T) {
	api := hcloudmock.New()
	srv := &hetzner.Server{ID: 1, Name: "web", ServerType: &hetzner.ServerType{Name: "cpx11"}}
	api.AddServer(srv)

	if err := Rescale(context.Background(), api, srv, "cpx21"); err != nil {
		t.Fatalf("Rescale: %v", err)
	}
	if srv.ServerType.Name != "cpx21" {
		t.Fatalf("server type after rescale = %q, want cpx21", srv.ServerType.Name)
	}
}

func TestRescaleSkipsIfAlreadyAtTarget(t *testing.T) {
	api := hcloudmock.New()
	srv := &hetzner.Server{ID: 1, Name: "web", ServerType: &hetzner.ServerType{Name: "cpx21"}}
	api.AddServer(srv)

	if err := Rescale(context.Background(), api, srv, "cpx21"); err != nil {
		t.Fatalf("Rescale (noop): %v", err)
	}
}

func TestRescaleReturnsUnavailableWhenTargetOutOfStock(t *testing.T) {
	api := hcloudmock.New()
	api.MarkUnavailable("cpx21")
	srv := &hetzner.Server{ID: 1, Name: "web", ServerType: &hetzner.ServerType{Name: "cpx11"}}
	api.AddServer(srv)

	err := Rescale(context.Background(), api, srv, "cpx21")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !hetzner.IsUnavailable(err) {
		t.Fatalf("err = %v, want unavailable", err)
	}
	_ = time.Second
}

func TestRescale_InvokesPhaseHook(t *testing.T) {
	api := hcloudmock.New()
	srv := &hetzner.Server{ID: 1, Name: "web", Status: hcloud.ServerStatusRunning, ServerType: &hetzner.ServerType{Name: "cpx11"}}
	api.AddServer(srv)

	var phases []string
	hook := func(p string) { phases = append(phases, p) }
	if err := RescaleWithHook(context.Background(), api, srv, "cpx21", hook); err != nil {
		t.Fatalf("Rescale: %v", err)
	}
	want := []string{"shutting_down", "changing_type", "powering_on", "done"}
	if len(phases) != len(want) {
		t.Fatalf("phases = %v, want %v", phases, want)
	}
	for i, p := range want {
		if phases[i] != p {
			t.Fatalf("phase[%d] = %q, want %q", i, phases[i], p)
		}
	}
}