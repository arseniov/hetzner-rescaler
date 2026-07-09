package rescaler

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/jonamat/hetzner-rescaler/internal/hcloudmock"
	"github.com/jonamat/hetzner-rescaler/internal/hetzner"
)

func TestFallbackSucceedsOnFirstTarget(t *testing.T) {
	api := hcloudmock.New()
	srv := &hetzner.Server{ID: 1, Name: "web", ServerType: &hetzner.ServerType{Name: "cpx11"}}
	api.AddServer(srv)

	used, err := RescaleWithFallback(context.Background(), api, srv, "cpx31", []string{"cpx31", "cpx21", "cpx11"})
	if err != nil {
		t.Fatalf("RescaleWithFallback: %v", err)
	}
	if used != "cpx31" {
		t.Fatalf("used = %q, want cpx31", used)
	}
}

func TestFallbackTriesNextOnUnavailable(t *testing.T) {
	api := hcloudmock.New()
	api.MarkUnavailable("cpx31")
	api.MarkUnavailable("cpx21")
	srv := &hetzner.Server{ID: 1, Name: "web", ServerType: &hetzner.ServerType{Name: "cpx11"}}
	api.AddServer(srv)

	used, err := RescaleWithFallback(context.Background(), api, srv, "cpx31", []string{"cpx31", "cpx21", "cpx11"})
	if err != nil {
		t.Fatalf("RescaleWithFallback: %v", err)
	}
	if used != "cpx11" {
		t.Fatalf("used = %q, want cpx11", used)
	}
}

func TestFallbackFailsWhenAllUnavailable(t *testing.T) {
	api := hcloudmock.New()
	api.MarkUnavailable("cpx31")
	api.MarkUnavailable("cpx21")
	api.MarkUnavailable("cpx11")
	// Server's current type is "cx11" — outside the chain — so the chain is
	// actually walked and Rescale's no-op short-circuit cannot rescue us.
	srv := &hetzner.Server{ID: 1, Name: "web", ServerType: &hetzner.ServerType{Name: "cx11"}}
	api.AddServer(srv)

	_, err := RescaleWithFallback(context.Background(), api, srv, "cpx31", []string{"cpx31", "cpx21", "cpx11"})
	if !errors.Is(err, ErrAllUnavailable) {
		t.Fatalf("err = %v, want ErrAllUnavailable", err)
	}
}

func TestFallbackStopsOnNonUnavailableError(t *testing.T) {
	api := hcloudmock.New()
	api.SetChangeTypeOverride(func(target *hetzner.ServerType) error {
		return errors.New("network blip")
	})
	srv := &hetzner.Server{ID: 1, Name: "web", ServerType: &hetzner.ServerType{Name: "cpx11"}}
	api.AddServer(srv)

	_, err := RescaleWithFallback(context.Background(), api, srv, "cpx21", []string{"cpx21", "cpx11"})
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if !strings.Contains(err.Error(), "network blip") {
		t.Fatalf("err = %v, want it to mention 'network blip'", err)
	}
}

// When the server is running and the fallback chain is walked, the server
// must be shut down exactly once — not once per failed chain entry.
// Otherwise the fallback chain pays N x (30s provisioner sleep) before the
// final successful entry, and a real Hetzner call would error because the
// server is already off.
func TestFallback_ShutdownHappensOnceWhenServerRunning(t *testing.T) {
	api := hcloudmock.New()
	api.MarkUnavailable("cpx31")
	api.MarkUnavailable("cpx21")
	srv := &hetzner.Server{
		ID:         1,
		Name:       "web",
		Status:     hcloud.ServerStatusRunning,
		ServerType: &hetzner.ServerType{Name: "cx11"},
	}
	api.AddServer(srv)

	used, err := RescaleWithFallback(context.Background(), api, srv, "cpx31", []string{"cpx31", "cpx21", "cpx11"})
	if err != nil {
		t.Fatalf("RescaleWithFallback: %v", err)
	}
	if used != "cpx11" {
		t.Fatalf("used = %q, want cpx11", used)
	}
	if got := api.ShutdownCount(); got != 1 {
		t.Fatalf("ShutdownCount = %d, want 1 (each chain entry must not re-shutdown an already-off server)", got)
	}
}