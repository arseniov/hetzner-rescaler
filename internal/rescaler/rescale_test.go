package rescaler

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
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

// CRITICAL contract: a rescale attempt MUST always leave the server
// running, regardless of whether the target type was available, the
// change_type action failed mid-flight, or a fallback was used. An
// off server is unreachable — silently taking production down — so
// "failed rescale" must never be conflated with "server stayed off".
func TestRescale_RestartsServerAfterUnavailable(t *testing.T) {
	api := hcloudmock.New()
	api.MarkUnavailable("cpx21")
	srv := &hetzner.Server{
		ID:         1,
		Name:       "web",
		Status:     hcloud.ServerStatusRunning,
		ServerType: &hetzner.ServerType{Name: "cpx11"},
	}
	api.AddServer(srv)

	err := Rescale(context.Background(), api, srv, "cpx21")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !hetzner.IsUnavailable(err) {
		t.Fatalf("err = %v, want unavailable", err)
	}
	if srv.Status != hcloud.ServerStatusRunning {
		t.Fatalf("srv.Status = %q, want running (server must always be restarted)", srv.Status)
	}
}

// Same contract, but for the fallback-chain path: when every entry in
// the chain is unavailable, the server must still be running when the
// chain gives up. This was the worst pre-fix case — three shut-downs
// and zero power-ons — because each RescaleWithHook invocation saw the
// stale "Running" status.
func TestRescale_RestartsServerAfterFallbackExhaustion(t *testing.T) {
	api := hcloudmock.New()
	api.MarkUnavailable("cpx31")
	api.MarkUnavailable("cpx21")
	api.MarkUnavailable("cpx11")
	srv := &hetzner.Server{
		ID:         1,
		Name:       "web",
		Status:     hcloud.ServerStatusRunning,
		ServerType: &hetzner.ServerType{Name: "cx11"},
	}
	api.AddServer(srv)

	_, err := RescaleWithFallback(context.Background(), api, srv, "cpx31", []string{"cpx31", "cpx21", "cpx11"})
	if !errors.Is(err, ErrAllUnavailable) {
		t.Fatalf("err = %v, want ErrAllUnavailable", err)
	}
	if srv.Status != hcloud.ServerStatusRunning {
		t.Fatalf("srv.Status = %q, want running (server must always be restarted)", srv.Status)
	}
}