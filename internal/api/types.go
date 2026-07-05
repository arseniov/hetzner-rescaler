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

// ServerResponse mirrors the store.Server struct with a pointer to
// string for PromoteState so null/omitted are distinguishable.
type ServerResponse struct {
	ID             int64    `json:"id"`
	ProjectID      int64    `json:"project_id"`
	HCloudServerID int      `json:"hcloud_server_id"`
	Name           string   `json:"name"`
	Label          string   `json:"label"`
	BaseServerType string   `json:"base_server_type"`
	TopServerType  string   `json:"top_server_type"`
	FallbackChain  []string `json:"fallback_chain"`
	Mode           string   `json:"mode"`
	PromoteState   *string  `json:"promote_state,omitempty"`
	Timezone       string   `json:"timezone"`
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
// UI's server-type picker. Available is the live availability flag.
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

// createProjectResponse is what POST /api/projects returns. It embeds
// ProjectResponse plus the auto-populated server tallies (from the
// initial sync) and any LastError from the fetch step (so the UI can
// surface a bad token without losing the created project row).
type createProjectResponse struct {
	ProjectResponse
	Added   []ServerResponse `json:"added"`
	Skipped []ServerResponse `json:"skipped"`
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