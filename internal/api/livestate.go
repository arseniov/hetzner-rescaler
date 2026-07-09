package api

import (
	"context"

	"github.com/jonamat/hetzner-rescaler/internal/store"
)

// LiveServerState is the slice of Hetzner server state the API exposes
// alongside the stored configuration. All fields are zero values when
// Hetzner is unreachable, the server has been deleted from Hetzner, or
// the project's API token is invalid — there is no failure path that
// turns into a non-2xx response from /api/servers.
type LiveServerState struct {
	Status      string
	CurrentType string
	Location    string
}

// liveServerState fetches the live state for a single server. It is
// the smallest unit the list/get handlers reuse; the list handler
// fans these out in parallel via liveStateMap.
//
// Soft-fail semantics: any error (token problems, network blip,
// Hetzner returning nil for a deleted server, etc.) returns the zero
// value. The handler decides whether to surface the omission to the
// client — for /api/servers the JSON tag omitempty drops the fields
// entirely so the web can render its own fallback.
//
// This function never panics on a nil store.Server: callers are
// expected to guard, but we double-check defensively because the
// caller in handleListServers is over a slice and a programming error
// there shouldn't take the whole endpoint down.
func (d Deps) liveServerState(ctx context.Context, srv *store.Server) LiveServerState {
	if srv == nil {
		return LiveServerState{}
	}
	if d.APIFor == nil {
		return LiveServerState{}
	}
	api, err := d.APIFor(srv.ProjectID)
	if err != nil {
		return LiveServerState{}
	}
	hs, err := api.GetServer(ctx, srv.HCloudServerID)
	if err != nil || hs == nil {
		return LiveServerState{}
	}
	out := LiveServerState{
		Status: string(hs.Status),
	}
	if hs.ServerType != nil {
		out.CurrentType = hs.ServerType.Name
	}
	if hs.Datacenter != nil && hs.Datacenter.Location != nil {
		out.Location = hs.Datacenter.Location.Name
	}
	return out
}