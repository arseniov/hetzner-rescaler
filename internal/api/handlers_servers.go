package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jonamat/hetzner-rescaler/internal/store"
)

func (d Deps) handleListServers(w http.ResponseWriter, r *http.Request) {
	servers, err := d.Store.ListAllServers()
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	out := make([]ServerResponse, 0, len(servers))
	for _, s := range servers {
		out = append(out, serverToResponse(s))
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
	writeJSON(w, http.StatusOK, serverToResponse(srv))
}

func (d Deps) handleCreateServer(w http.ResponseWriter, r *http.Request) {
	var req CreateServerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if req.ProjectID == 0 || req.HCloudServerID == 0 || req.Name == "" ||
		req.BaseServerType == "" || req.TopServerType == "" ||
		len(req.FallbackChain) == 0 || req.Mode == "" || req.Timezone == "" {
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
	writeJSON(w, http.StatusCreated, serverToResponse(srv))
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
		len(req.FallbackChain) == 0 || req.Mode == "" || req.Timezone == "" {
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
	writeJSON(w, http.StatusOK, serverToResponse(existing))
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

// serverToResponse converts a *store.Server to its API projection.
// Canonical definition lives here (Task 5); the previous duplicate in
// handlers_projects.go was removed in favor of this one.
func serverToResponse(s *store.Server) ServerResponse {
	if s == nil {
		return ServerResponse{}
	}
	return ServerResponse{
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
	}
}