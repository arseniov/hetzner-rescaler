package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"sync"

	"github.com/jonamat/hetzner-rescaler/internal/store"
)

// liveStateFanOutLimit caps the number of parallel Hetzner GetServer
// calls in handleListServers. Eight keeps the worst-case request
// latency well under one second even when the Hetzner API is slow,
// without hammering their rate limiter on a large fleet.
const liveStateFanOutLimit = 8

func (d Deps) handleListServers(w http.ResponseWriter, r *http.Request) {
	servers, err := d.Store.ListAllServers()
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	live := liveStateMap(r.Context(), d, servers)

	// Pending events are read serially per row (one indexed lookup each).
	// Acceptable here: matches the previous inline pattern, and SQLite is
	// already on the critical path for the rest of the request.
	pendingMap := make(map[int64]*store.Event, len(servers))
	for _, s := range servers {
		if e, err := d.Store.ActivePendingEventForServer(s.ID); err == nil {
			pendingMap[s.ID] = e
		}
	}

	out := make([]ServerResponse, 0, len(servers))
	for _, s := range servers {
		out = append(out, serverToResponse(s, live[s.ID], pendingMap[s.ID]))
	}
	writeJSON(w, http.StatusOK, out)
}

func (d Deps) handleGetServer(w http.ResponseWriter, r *http.Request) {
	id, ok := pathInt64(r, "id")
	if !ok {
		writeJSONError(w, http.StatusBadRequest, "invalid id")
		return
	}
	srv, err := d.Store.GetServer(id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeJSONError(w, http.StatusNotFound, "server not found")
			return
		}
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	live := d.liveServerState(r.Context(), srv)
	pending, _ := d.Store.ActivePendingEventForServer(srv.ID)
	writeJSON(w, http.StatusOK, serverToResponse(srv, live, pending))
}

func (d Deps) handleCreateServer(w http.ResponseWriter, r *http.Request) {
	var req CreateServerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if req.ProjectID == 0 || req.HCloudServerID == 0 || req.Name == "" ||
		req.BaseServerType == "" || req.TopServerType == "" ||
		req.Mode == "" || req.Timezone == "" {
		writeJSONError(w, http.StatusBadRequest, "missing required fields")
		return
	}
	if _, err := d.Store.GetProject(req.ProjectID); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeJSONError(w, http.StatusBadRequest, "unknown project_id")
			return
		}
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	srv, err := d.Store.CreateServer(req.ProjectID, store.Server{
		HCloudServerID: req.HCloudServerID,
		Name:           req.Name,
		Label:          req.Label,
		BaseServerType: req.BaseServerType,
		TopServerType:  req.TopServerType,
		FallbackChain:  req.FallbackChain,
		Mode:           req.Mode,
		Timezone:       req.Timezone,
	})
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	live := d.liveServerState(r.Context(), srv)
	writeJSON(w, http.StatusCreated, serverToResponse(srv, live, nil))
}

func (d Deps) handleUpdateServer(w http.ResponseWriter, r *http.Request) {
	id, ok := pathInt64(r, "id")
	if !ok {
		writeJSONError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var req UpdateServerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if req.Name == "" || req.BaseServerType == "" || req.TopServerType == "" ||
		req.Mode == "" || req.Timezone == "" {
		writeJSONError(w, http.StatusBadRequest, "missing required fields")
		return
	}
	existing, err := d.Store.GetServer(id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeJSONError(w, http.StatusNotFound, "server not found")
			return
		}
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	existing.Name = req.Name
	existing.Label = req.Label
	existing.BaseServerType = req.BaseServerType
	existing.TopServerType = req.TopServerType
	existing.FallbackChain = req.FallbackChain
	existing.Mode = req.Mode
	existing.Timezone = req.Timezone
	if err := d.Store.UpdateServer(*existing); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	live := d.liveServerState(r.Context(), existing)
	writeJSON(w, http.StatusOK, serverToResponse(existing, live, nil))
}

func (d Deps) handleDeleteServer(w http.ResponseWriter, r *http.Request) {
	id, ok := pathInt64(r, "id")
	if !ok {
		writeJSONError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if _, err := d.Store.GetServer(id); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeJSONError(w, http.StatusNotFound, "server not found")
			return
		}
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if err := d.Store.DeleteServer(id); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// liveStateMap fetches live state for every server in parallel, capped
// at liveStateFanOutLimit concurrent calls. Each call soft-fails on
// its own (returning a zero LiveServerState), so the worst case is
// "all servers come back without status/current_type" — never an
// error response. Map keys are store.Server.ID.
func liveStateMap(ctx context.Context, d Deps, servers []*store.Server) map[int64]LiveServerState {
	out := make(map[int64]LiveServerState, len(servers))
	if len(servers) == 0 || d.APIFor == nil {
		return out
	}

	type result struct {
		id   int64
		live LiveServerState
	}
	results := make(chan result, len(servers))
	sem := make(chan struct{}, liveStateFanOutLimit)

	var wg sync.WaitGroup
	for _, s := range servers {
		if s == nil {
			continue
		}
		wg.Add(1)
		sem <- struct{}{}
		go func(srv *store.Server) {
			defer wg.Done()
			defer func() { <-sem }()
			// Re-check nil in the goroutine — defensive, matches the
			// single-call path in liveServerState.
			if srv == nil {
				results <- result{}
				return
			}
			live := d.liveServerState(ctx, srv)
			results <- result{id: srv.ID, live: live}
		}(s)
	}
	wg.Wait()
	close(results)

	for r := range results {
		if r.id != 0 {
			out[r.id] = r.live
		}
	}
	return out
}

// serverToResponse converts a *store.Server to its API projection,
// merging in the live Hetzner state when present. A zero LiveServerState
// produces omitempty-absent fields, so callers that skip the live
// fetch (or hit a Hetzner failure) still get a well-formed JSON body.
// A nil pending leaves pending_event absent from the response.
func serverToResponse(s *store.Server, live LiveServerState, pending *store.Event) ServerResponse {
	if s == nil {
		return ServerResponse{}
	}
	resp := ServerResponse{
		ID:             s.ID,
		ProjectID:      s.ProjectID,
		HCloudServerID: s.HCloudServerID,
		Name:           s.Name,
		Label:          s.Label,
		BaseServerType: s.BaseServerType,
		TopServerType:  s.TopServerType,
		FallbackChain:  s.FallbackChain,
		Mode:           s.Mode,
		PromoteState:   s.PromoteState,
		Timezone:       s.Timezone,
		Status:         live.Status,
		CurrentType:    live.CurrentType,
		CreatedAt:      s.CreatedAt,
		UpdatedAt:      s.UpdatedAt,
	}
	if pending != nil {
		pe := eventToResponse(pending)
		resp.PendingEvent = &pe
	}
	return resp
}
