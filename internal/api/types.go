package api

import "time"

// --- Response types (what handlers return) ---

// ProjectResponse is what /api/projects returns. The token is NEVER
// included — only a boolean indicating that one is stored.
type ProjectResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	HasToken  bool      `json:"has_token"`
	LastError string    `json:"last_error"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ServerResponse is the API projection of a store.Server, plus live state
// fetched from Hetzner for every list/get request.
//
// Field sources:
//   - store-backed fields: id, project_id, hcloud_server_id, name, label,
//     base_server_type, top_server_type, fallback_chain, mode,
//     promote_state, timezone, created_at, updated_at, pending_event
//   - live Hetzner state: status (running/initializing/starting/stopping/
//     off/deleting), current_type (the live ServerType.Name as reported
//     by Hetzner right now)
//
// Live fields are omitempty: a Hetzner API failure leaves them absent
// from the response so the web can fall back to its own derived state
// instead of crashing the whole endpoint.
type ServerResponse struct {
	ID             int64     `json:"id"`
	ProjectID      int64     `json:"project_id"`
	HCloudServerID int       `json:"hcloud_server_id"`
	Name           string    `json:"name"`
	Label          string    `json:"label"`
	BaseServerType string    `json:"base_server_type"`
	TopServerType  string    `json:"top_server_type"`
	FallbackChain  []string  `json:"fallback_chain"`
	Mode           string    `json:"mode"`
	PromoteState   *string   `json:"promote_state,omitempty"`
	Timezone       string    `json:"timezone"`
	Status         string    `json:"status,omitempty"`
	CurrentType    string    `json:"current_type,omitempty"`
	Location       string    `json:"location,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	// PendingEvent is the in-flight rescale_pending event for this server,
	// if any. Omitted when the server is idle so the client can distinguish
	// "no rescale" from "rescale in progress" without a follow-up GET.
	PendingEvent *EventResponse `json:"pending_event,omitempty"`
}

// WindowResponse mirrors store.Window.
type WindowResponse struct {
	ID         int64  `json:"id"`
	ServerID   int64  `json:"server_id"`
	Label      string `json:"label"`
	DaysOfWeek int    `json:"days_of_week"`
	StartTime  string `json:"start_time"`
	StopTime   string `json:"stop_time"`
	TargetType string `json:"target_type"`
	Enabled    bool   `json:"enabled"`
}

// EventResponse mirrors store.Event. ok is rendered as a bool.
type EventResponse struct {
	ID          int64     `json:"id"`
	ServerID    int64     `json:"server_id"`
	Kind        string    `json:"kind"`
	FromType    string    `json:"from_type,omitempty"`
	ToType      string    `json:"to_type,omitempty"`
	StartedAt   time.Time `json:"started_at"`
	FinishedAt  time.Time `json:"finished_at"`
	OK          bool      `json:"ok"`
	Error       string    `json:"error,omitempty"`
	TriggeredBy string    `json:"triggered_by"`
}

// ServerTypeResponse is the projection of a Hetzner server type for the
// UI's server-type picker. All seven fields are populated on every
// non-nil source row: description, cores, memory/disk in GB, the live
// Available flag (derived from Pricings length), and the first
// pricing entry's monthly EUR gross parsed to float32.
type ServerTypeResponse struct {
	Name            string  `json:"name"`
	Description     string  `json:"description,omitempty"`
	Cores           int     `json:"cores"`
	MemoryGB        float32 `json:"memory_gb"`
	DiskGB          float32 `json:"disk_gb"`
	Available       bool    `json:"available"`
	PriceMonthlyEUR float32 `json:"price_monthly_eur,omitempty"`
}

// --- Request types (what handlers accept) ---

// CreateProjectRequest is the body for POST /api/projects.
type CreateProjectRequest struct {
	Name        string `json:"name"`
	HCloudToken string `json:"hcloud_token"`
}

// RefreshProjectResponse is what /api/projects/:id/refresh returns — a
// tally of newly-added servers and skipped ones.
type RefreshProjectResponse struct {
	Added   []ServerResponse `json:"added"`
	Skipped []ServerResponse `json:"skipped"`
}

// CreateProjectResponse is what POST /api/projects returns. It embeds
// ProjectResponse plus the auto-populated server tallies (from the
// initial sync) and a LastError from the fetch step so the UI can
// surface a bad token without losing the created project row.
//
// Exported because the web client types its createProject return as
// CreateProjectResult (which mirrors this shape). Kept in sync with
// the ProjectResponse embedded type so web/src/lib/types.ts sees
// `added`, `skipped`, `last_error`, plus every Project field flat on
// the response.
type CreateProjectResponse struct {
	ProjectResponse
	Added     []ServerResponse `json:"added"`
	Skipped   []ServerResponse `json:"skipped"`
	LastError string           `json:"last_error,omitempty"`
}

// CreateServerRequest is the body for POST /api/servers.
type CreateServerRequest struct {
	ProjectID      int64    `json:"project_id"`
	HCloudServerID int      `json:"hcloud_server_id"`
	Name           string   `json:"name"`
	Label          string   `json:"label"`
	BaseServerType string   `json:"base_server_type"`
	TopServerType  string   `json:"top_server_type"`
	FallbackChain  []string `json:"fallback_chain"`
	Mode           string   `json:"mode"`
	Timezone       string   `json:"timezone"`
}

// UpdateServerRequest is the body for PUT /api/servers/:id.
type UpdateServerRequest struct {
	Name           string   `json:"name"`
	Label          string   `json:"label"`
	BaseServerType string   `json:"base_server_type"`
	TopServerType  string   `json:"top_server_type"`
	FallbackChain  []string `json:"fallback_chain"`
	Mode           string   `json:"mode"`
	Timezone       string   `json:"timezone"`
}

// CreateWindowRequest is the body for POST /api/servers/:id/windows.
type CreateWindowRequest struct {
	Label      string `json:"label"`
	DaysOfWeek int    `json:"days_of_week"`
	StartTime  string `json:"start_time"`
	StopTime   string `json:"stop_time"`
	TargetType string `json:"target_type"`
	Enabled    bool   `json:"enabled"`
}

// UpdateWindowRequest is the body for PUT /api/windows/:wid.
type UpdateWindowRequest = CreateWindowRequest

// ConfirmRequest is the body for any destructive action that requires
// an explicit confirmation flag.
type ConfirmRequest struct {
	Confirm bool `json:"confirm"`
}

// RescaleRequest is the body for POST /api/servers/:id/rescale.
// Direction is "up" (→ top_server_type) or "down" (→ base_server_type).
type RescaleRequest struct {
	Direction string `json:"direction"`
	Confirm   bool   `json:"confirm"`
}

// --- Helpers for handlers ---

// toProjectResponse converts a *store.Project to its API projection.
// (Defined in handlers_projects.go when that task lands; declared here
// so types.go is self-contained for review.)

// toServerResponse converts a *store.Server to its API projection.
// (Defined in handlers_servers.go.)

// toWindowResponse converts a *store.Window to its API projection.
// (Defined in handlers_windows.go.)