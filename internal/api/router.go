// Package api exposes the engine's projects, servers, windows, events, and
// rescale actions over HTTP. It is intentionally thin: handlers translate
// JSON ↔ store types and call into store, rescaler, and scheduler.
package api

import (
	"encoding/json"
	"net/http"

	"github.com/jonamat/hetzner-rescaler/internal/crypto"
	"github.com/jonamat/hetzner-rescaler/internal/hetzner"
	"github.com/jonamat/hetzner-rescaler/internal/rescaler"
	"github.com/jonamat/hetzner-rescaler/internal/store"
)

// Deps holds dependencies the router needs. It is passed to NewRouter so
// tests can construct a router without touching the filesystem.
type Deps struct {
	// InternalToken is the shared secret expected in the X-Internal-Token
	// header on all /api/* calls except /api/healthz. Required.
	InternalToken string

	// SessionSecret is the same Better Auth secret that rescaler-web
	// uses to sign session cookies (BETTER_AUTH_SECRET). When non-empty
	// the auth middleware also admits requests carrying a valid
	// better-auth.session_token cookie. Empty disables the cookie path —
	// every authenticated request must then carry X-Internal-Token.
	SessionSecret string

	// Store is the SQLite-backed persistence layer. Required for handlers
	// that read or mutate engine state, and for the session lookup
	// performed by RequireAuth.
	Store *store.Store

	// Keyring is the AES-256 key used to seal Hetzner tokens before
	// persistence. If nil the handler falls back to KeyFromEnv(), which
	// reads RESCALER_ENCRYPTION_KEY from the environment (and generates
	// a fresh key if unset).
	Keyring *crypto.Keyring

	// APIFor returns a hetzner.API for a given project ID. Required for
	// handlers that talk to Hetzner (refresh, server-types).
	APIFor func(projectID int64) (hetzner.API, error)

	// Manager runs rescale work asynchronously. handleRescale calls
	// Manager.Submit and returns 202 immediately with the pending event
	// ID; Manager.runRescale walks the rescale phases in a goroutine
	// and writes terminal event rows.
	Manager *rescaler.Manager
}

// NewRouter builds the HTTP mux. /api/healthz is always registered.
// Every other /api/* route is registered as it is added.
func NewRouter(deps Deps) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	// All /api/* routes (except /api/healthz) require either:
	//   - a valid X-Internal-Token header (CLI scripts), or
	//   - a verified Better Auth session cookie (browser SPA).
	auth := RequireAuth(deps.InternalToken, deps.SessionSecret, deps.Store)
	streamAuth := eventsStreamAuth(deps.InternalToken, deps.SessionSecret, deps.Store, http.HandlerFunc(deps.handleEventsStream))

	// Project routes
	mux.Handle("GET /api/projects", auth(http.HandlerFunc(deps.handleListProjects)))
	mux.Handle("POST /api/projects", auth(http.HandlerFunc(deps.handleCreateProject)))
	mux.Handle("DELETE /api/projects/{id}", auth(http.HandlerFunc(deps.handleDeleteProject)))
	mux.Handle("POST /api/projects/{id}/refresh", auth(http.HandlerFunc(deps.handleRefreshProject)))

	// Server routes
	mux.Handle("GET /api/servers", auth(http.HandlerFunc(deps.handleListServers)))
	mux.Handle("GET /api/servers/{id}", auth(http.HandlerFunc(deps.handleGetServer)))
	mux.Handle("POST /api/servers", auth(http.HandlerFunc(deps.handleCreateServer)))
	mux.Handle("PUT /api/servers/{id}", auth(http.HandlerFunc(deps.handleUpdateServer)))
	mux.Handle("DELETE /api/servers/{id}", auth(http.HandlerFunc(deps.handleDeleteServer)))

	// Window routes
	mux.Handle("GET /api/servers/{id}/windows", auth(http.HandlerFunc(deps.handleListWindows)))
	mux.Handle("POST /api/servers/{id}/windows", auth(http.HandlerFunc(deps.handleCreateWindow)))
	mux.Handle("PUT /api/windows/{wid}", auth(http.HandlerFunc(deps.handleUpdateWindow)))
	mux.Handle("DELETE /api/windows/{wid}", auth(http.HandlerFunc(deps.handleDeleteWindow)))

	// Action routes
	mux.Handle("POST /api/servers/{id}/rescale", auth(http.HandlerFunc(deps.handleRescale)))
	mux.Handle("POST /api/servers/{id}/promote", auth(http.HandlerFunc(deps.handlePromote)))
	mux.Handle("POST /api/servers/{id}/demote", auth(http.HandlerFunc(deps.handleDemote)))

	// Event routes
	mux.Handle("GET /api/servers/{id}/events", auth(http.HandlerFunc(deps.handleServerEvents)))
	mux.Handle("GET /api/events", auth(http.HandlerFunc(deps.handleGlobalEvents)))
	mux.Handle("GET /api/events/stream", streamAuth)

	// Server-type route (proxies to Hetzner via the first project's API)
	mux.Handle("GET /api/server-types", auth(http.HandlerFunc(deps.handleServerTypes)))

	// Metrics route (dashboard aggregations)
	mux.Handle("GET /api/metrics", auth(http.HandlerFunc(deps.handleMetrics)))

	return mux
}
