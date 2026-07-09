package cmd

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jonamat/hetzner-rescaler/internal/api"
	"github.com/jonamat/hetzner-rescaler/internal/broadcast"
	"github.com/jonamat/hetzner-rescaler/internal/crypto"
	"github.com/jonamat/hetzner-rescaler/internal/hetzner"
	"github.com/jonamat/hetzner-rescaler/internal/rescaler"
	"github.com/jonamat/hetzner-rescaler/internal/scheduler"
	"github.com/jonamat/hetzner-rescaler/internal/store"

	"github.com/spf13/cobra"
)

// keyringHolder is set by runServe so the API package's apiFactory can
// decrypt per-project Hetzner tokens. It is package-scoped so apiFactory
// (also in this package) can read it without threading it through every
// handler.
var keyringHolder *crypto.Keyring

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run the HTTP API server + static SPA + scheduler (phase 2)",
	Long: `Runs the rescaler in 'serve' mode: a loopback HTTP server that
exposes the /api/* JSON API and serves the embedded SvelteKit SPA at all
other paths. The scheduler runs in the same process so rescale events
fire whether the trigger came from the CLI, the API, or a scheduled window.

Required environment:
  RESCALER_INTERNAL_TOKEN  shared secret between the SPA and the API

Optional environment:
  RESCALER_HTTP_ADDR       listen address (default 127.0.0.1:8080)
  RESCALER_DB_PATH         SQLite file (default ~/.hetzner-rescaler/db.sqlite)
  RESCALER_TOKEN_ENCRYPTION_KEY  hex-encoded 32-byte AES-GCM key`,
	RunE: runServe,
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

func runServe(cmd *cobra.Command, args []string) error {
	token := os.Getenv("RESCALER_INTERNAL_TOKEN")
	if token == "" {
		return fmt.Errorf("serve: RESCALER_INTERNAL_TOKEN is required (the SPA needs it to call /api/*)")
	}

	// BETTER_AUTH_SECRET is optional. When set (recommended), the API
	// middleware will additionally admit requests that carry a valid
	// better-auth.session_token cookie, giving the SPA a second auth
	// path that doesn't depend on the internal token. When unset, the
	// API reverts to InternalToken-only (matches CLI/scripts).
	sessionSecret := os.Getenv("BETTER_AUTH_SECRET")

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	st, key, err := openStoreAndKeyring()
	if err != nil {
		return fmt.Errorf("serve: %w", err)
	}
	defer st.Close()
	keyringHolder = key

	// Live events: broadcast every inserted event row to in-process subscribers
	// (the SSE handler added in a later task).
	eventHub := broadcast.NewHub[store.Event]()
	st.SetBroadcastHub(eventHub)

	// Server-lifecycle hub: receive every CreateServer / UpdateServer /
	// DeleteServer write so the scheduler can add or remove per-server
	// goroutines without polling.
	lifecycleHub := broadcast.NewHub[store.ServerLifecycleEvent]()
	st.SetServerLifecycleHub(lifecycleHub)

	// Scheduler: per-server goroutines that drive scheduled / auto_promote
	// rescales. The apiResolve closure uses the same per-project apiFactory
	// the rescaler manager uses, so multi-project deployments schedule
	// every project's servers.
	sched := scheduler.New(st, func(_ context.Context, projectID int64) (hetzner.API, error) {
		return apiFactory(st)(projectID)
	}, scheduler.RealClock{}, 30*time.Second)

	// Rescale manager: owns the async rescale goroutines. Built first so it
	// can Start (recovery of orphaned pending rows) before the API begins
	// serving. The resolver wires Manager to the per-project API factory
	// so each Submit gets the correct project's hetzner.API.
	rescalerManager := rescalerManager(st)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	if err := rescalerManager.Start(ctx); err != nil {
		return fmt.Errorf("serve: rescaler manager start: %w", err)
	}
	defer func() {
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()
		_ = rescalerManager.Shutdown(shutdownCtx)
	}()

	go func() {
		if err := sched.Attach(ctx, lifecycleHub); err != nil && !errors.Is(err, context.Canceled) {
			logger.Error("scheduler attach", "err", err)
		}
	}()
	defer sched.Stop()

	deps := api.Deps{
		InternalToken: token,
		SessionSecret: sessionSecret,
		Store:         st,
		Keyring:       key,
		APIFor:        apiFactory(st),
		Manager:       rescalerManager,
	}

	// Compose: the API server only owns /api/*. The SPA is served by a
	// separate SvelteKit adapter-node process (rescaler-web); a reverse
	// proxy (Caddy) routes /api/auth/* → rescaler-web and the remaining
	// /api/* → rescaler-api.
	apiMux := api.NewRouter(deps)
	handler := apiMux

	addr := os.Getenv("RESCALER_HTTP_ADDR")
	if addr == "" {
		addr = "127.0.0.1:8080"
	}
	srv := &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: 10 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		logger.Info("serve: listening", "addr", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
		close(errCh)
	}()
	select {
	case <-ctx.Done():
	case err := <-errCh:
		if err != nil {
			return fmt.Errorf("serve: ListenAndServe: %w", err)
		}
	}
	logger.Info("serve: shutting down")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	return srv.Shutdown(shutdownCtx)
}

// apiFactory returns the Deps.APIFor function: it looks up the project's
// encrypted token in the store, decrypts it with the package-scoped
// keyring, and constructs a hetzner.API.
func apiFactory(s *store.Store) func(projectID int64) (hetzner.API, error) {
	return func(projectID int64) (hetzner.API, error) {
		p, err := s.GetProject(projectID)
		if err != nil {
			return nil, fmt.Errorf("get project: %w", err)
		}
		tok, err := keyringHolder.Open(p.HCloudTokenEncrypted, p.HCloudTokenNonce)
		if err != nil {
			return nil, fmt.Errorf("decrypt token: %w", err)
		}
		return hetzner.NewClient(string(tok))
	}
}

// rescalerManager builds the rescaler.Manager wired to the project's
// per-project API factory. Start/Stop are the caller's responsibility
// (runServe calls them).
func rescalerManager(s *store.Store) *rescaler.Manager {
	m := rescaler.NewManager(s)
	factory := apiFactory(s)
	m.SetAPIResolver(func(ctx context.Context, projectID int64) (hetzner.API, error) {
		return factory(projectID)
	})
	return m
}