// Package api exposes the engine's projects, servers, windows, events, and
// rescale actions over HTTP. It is intentionally thin: handlers translate
// JSON ↔ store types and call into store, rescaler, and scheduler.
package api

import (
	"encoding/json"
	"net/http"
)

// Deps holds dependencies the router needs. It is passed to NewRouter so
// tests can construct a router without touching the filesystem.
type Deps struct {
	// InternalToken is the shared secret expected in the X-Internal-Token
	// header on all /api/* calls except /api/healthz. Required.
	InternalToken string

	// Store is the SQLite-backed persistence layer. Required for handlers
	// that read or mutate engine state. Typed as `any` in this stub; the
	// concrete type is wired in Task 4.
	Store any

	// Hetzner factories: a function that returns a hetzner.API for a
	// given project ID. Nil until Task 4 (project handlers) wires it up.
	APIFor func(projectID int64) (HetznerAPI, error)
}

// HetznerAPI is the minimal Hetzner surface the API package needs.
// The full interface lives in internal/hetzner; this is declared here so
// router.go doesn't depend on that package's transitive deps.
type HetznerAPI interface {
	ListServers(ctx ctxLike) ([]HetznerServer, error)
	GetServer(ctx ctxLike, id int) (*HetznerServer, error)
	ListServerTypes(ctx ctxLike) ([]HetznerServerType, error)
}

// ctxLike avoids importing "context" in this stub interface.
// The real interface in Task 5+ uses context.Context directly.
type ctxLike interface {
	Deadline() (deadline interface{}, ok bool)
}

// HetznerServer / HetznerServerType are minimal projections; the real
// types come from internal/hetzner in later tasks.
type (
	HetznerServer     struct{ ID int; Name string; ServerType string }
	HetznerServerType struct{ Name string; Available bool }
)

// NewRouter builds the HTTP mux. /api/healthz is always registered.
// Every other /api/* route is registered as it is added in later tasks.
func NewRouter(deps Deps) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})
	return mux
}
