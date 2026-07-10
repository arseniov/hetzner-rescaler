package rescaler

import (
	"context"
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/jonamat/hetzner-rescaler/internal/hcloudmock"
	"github.com/jonamat/hetzner-rescaler/internal/hetzner"
	"github.com/jonamat/hetzner-rescaler/internal/store"
)

func newTestStore(t *testing.T) *store.Store {
	t.Helper()
	st, err := store.Open(filepath.Join(t.TempDir(), "avail.db"))
	if err != nil {
		t.Fatalf("store.Open: %v", err)
	}
	t.Cleanup(func() { st.Close() })
	p, err := st.CreateProject("p", []byte("tok"), []byte("nonce12byts"))
	if err != nil {
		t.Fatalf("CreateProject: %v", err)
	}
	_ = p
	return st
}

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

func TestFallback_UnavailableEntrySkippedWithoutShutdown(t *testing.T) {
	api := hcloudmock.New()
	api.SetLocations("cpx11", "fsn1", false) // cpx11 unavailable in fsn1
	api.SetLocations("cpx21", "fsn1", true)  // cpx21 available — should be picked
	st := newTestStore(t)
	srv := &hetzner.Server{
		ID:         1,
		Name:       "web",
		Status:     hcloud.ServerStatusRunning,
		ServerType: &hetzner.ServerType{Name: "cx11"},
		Datacenter: &hetzner.Datacenter{Location: &hetzner.Location{Name: "fsn1"}},
	}
	if _, err := st.CreateServer(1, store.Server{HCloudServerID: 1, Name: "w", BaseServerType: "cx11", TopServerType: "cpx21", FallbackChain: []string{"cpx11", "cpx21"}, Mode: "manual", Timezone: "UTC"}); err != nil {
		t.Fatalf("CreateServer: %v", err)
	}

	used, err := RescaleWithFallbackWithHook(context.Background(), api, srv, "cpx11", []string{"cpx11", "cpx21"}, st, 1, "operator", nil)
	if err != nil {
		t.Fatalf("err = %v, want nil (cpx21 should succeed)", err)
	}
	if used != "cpx21" {
		t.Fatalf("used = %q, want cpx21", used)
	}
	if got := api.ShutdownCount(); got != 1 {
		t.Fatalf("ShutdownCount = %d, want 1 (only cpx21 should trigger shutdown, not cpx11)", got)
	}

	events, _ := st.ListEventsByServer(1, 100)
	skipped := 0
	for _, e := range events {
		switch e.Kind {
		case "rescale_skipped":
			skipped++
			if e.FromType != "cpx11" {
				t.Fatalf("skipped event FromType = %q, want cpx11", e.FromType)
			}
			if !strings.Contains(e.Error, "unavailable in fsn1") {
				t.Fatalf("skipped event Error = %q, want it to contain \"unavailable in fsn1\"", e.Error)
			}
			if e.OK {
				t.Fatalf("skipped event OK = true, want false")
			}
			if e.TriggeredBy != "operator" {
				t.Fatalf("skipped event TriggeredBy = %q, want operator", e.TriggeredBy)
			}
			if e.Phase != "pre_check" {
				t.Fatalf("skipped event Phase = %q, want pre_check", e.Phase)
			}
		}
	}
	if skipped != 1 {
		t.Fatalf("rescale_skipped count = %d, want 1", skipped)
	}
}

func TestFallback_AllUnavailableReturnsErrAllUnavailable(t *testing.T) {
	api := hcloudmock.New()
	api.SetLocations("cpx11", "fsn1", false)
	api.SetLocations("cpx21", "fsn1", false)
	st := newTestStore(t)
	if _, err := st.CreateServer(1, store.Server{HCloudServerID: 1, Name: "w", BaseServerType: "cpx11", TopServerType: "cpx21", FallbackChain: []string{"cpx11", "cpx21"}, Mode: "manual", Timezone: "UTC"}); err != nil {
		t.Fatalf("CreateServer: %v", err)
	}
	srv := &hetzner.Server{
		ID:         1,
		Name:       "web",
		Status:     hcloud.ServerStatusRunning,
		ServerType: &hetzner.ServerType{Name: "cx11"},
		Datacenter: &hetzner.Datacenter{Location: &hetzner.Location{Name: "fsn1"}},
	}

	_, err := RescaleWithFallbackWithHook(context.Background(), api, srv, "cpx11", []string{"cpx11", "cpx21"}, st, 1, "operator", nil)
	if !errors.Is(err, ErrAllUnavailable) {
		t.Fatalf("err = %v, want ErrAllUnavailable", err)
	}
	if got := api.ShutdownCount(); got != 0 {
		t.Fatalf("ShutdownCount = %d, want 0 (no chain entry should have triggered a shutdown)", got)
	}
	events, _ := st.ListEventsByServer(1, 100)
	skipped := 0
	for _, e := range events {
		if e.Kind == "rescale_skipped" {
			skipped++
		}
	}
	if skipped != 2 {
		t.Fatalf("rescale_skipped count = %d, want 2", skipped)
	}
}

func TestFallback_PreCheckAPIErrorFailsOpen(t *testing.T) {
	api := hcloudmock.New()
	api.SetGetServerTypeError("cpx21", errors.New("transient"))
	st := newTestStore(t)
	if _, err := st.CreateServer(1, store.Server{HCloudServerID: 1, Name: "w", BaseServerType: "cpx11", TopServerType: "cpx21", FallbackChain: []string{"cpx21"}, Mode: "manual", Timezone: "UTC"}); err != nil {
		t.Fatalf("CreateServer: %v", err)
	}
	srv := &hetzner.Server{
		ID:         1,
		Name:       "web",
		Status:     hcloud.ServerStatusRunning,
		ServerType: &hetzner.ServerType{Name: "cx11"},
		Datacenter: &hetzner.Datacenter{Location: &hetzner.Location{Name: "fsn1"}},
	}

	used, err := RescaleWithFallbackWithHook(context.Background(), api, srv, "cpx21", []string{"cpx21"}, st, 1, "operator", nil)
	if err != nil {
		t.Fatalf("err = %v, want nil (fail-open: the API error should not block the rescale)", err)
	}
	if used != "cpx21" {
		t.Fatalf("used = %q, want cpx21", used)
	}
}

func TestFallback_NilStoreDoesNotEmit(t *testing.T) {
	api := hcloudmock.New()
	api.SetLocations("cpx11", "fsn1", false)
	srv := &hetzner.Server{
		ID:         1,
		Name:       "web",
		Status:     hcloud.ServerStatusRunning,
		ServerType: &hetzner.ServerType{Name: "cx11"},
		Datacenter: &hetzner.Datacenter{Location: &hetzner.Location{Name: "fsn1"}},
	}

	_, err := RescaleWithFallback(context.Background(), api, srv, "cpx11", []string{"cpx11"})
	if !errors.Is(err, ErrAllUnavailable) {
		t.Fatalf("err = %v, want ErrAllUnavailable", err)
	}
	if got := api.ShutdownCount(); got != 0 {
		t.Fatalf("ShutdownCount = %d, want 0", got)
	}
}