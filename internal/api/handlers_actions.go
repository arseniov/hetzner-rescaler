package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jonamat/hetzner-rescaler/internal/store"
)

func (d Deps) handleRescale(w http.ResponseWriter, r *http.Request) {
	srv, ok := d.resolveServer(w, r)
	if !ok {
		return
	}
	var req RescaleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if !req.Confirm {
		writeJSONError(w, http.StatusBadRequest, "confirm: true required")
		return
	}
	var target string
	switch req.Direction {
	case "up":
		target = srv.TopServerType
	case "down":
		target = srv.BaseServerType
	default:
		writeJSONError(w, http.StatusBadRequest, "direction must be 'up' or 'down'")
		return
	}
	if d.Rescaler == nil {
		writeJSONError(w, http.StatusServiceUnavailable, "rescaler not configured")
		return
	}
	if err := d.Rescaler(r.Context(), srv, target); err != nil {
		writeJSONError(w, http.StatusBadGateway, "rescale failed: "+err.Error())
		return
	}
	writeJSON(w, http.StatusAccepted, map[string]string{"status": "rescale initiated", "target": target})
}

func (d Deps) handlePromote(w http.ResponseWriter, r *http.Request) {
	srv, ok := d.resolveServer(w, r)
	if !ok {
		return
	}
	var req ConfirmRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || !req.Confirm {
		writeJSONError(w, http.StatusBadRequest, "confirm: true required")
		return
	}
	if srv.Mode != "auto_promote" {
		writeJSONError(w, http.StatusBadRequest, "promote is only valid in auto_promote mode")
		return
	}
	state := "promote_requested"
	srv.PromoteState = &state
	if err := d.Store.UpdateServer(*srv); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusAccepted, map[string]string{"status": "promote requested"})
}

func (d Deps) handleDemote(w http.ResponseWriter, r *http.Request) {
	srv, ok := d.resolveServer(w, r)
	if !ok {
		return
	}
	var req ConfirmRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || !req.Confirm {
		writeJSONError(w, http.StatusBadRequest, "confirm: true required")
		return
	}
	if srv.Mode != "auto_promote" {
		writeJSONError(w, http.StatusBadRequest, "demote is only valid in auto_promote mode")
		return
	}
	state := "demote_requested"
	srv.PromoteState = &state
	if err := d.Store.UpdateServer(*srv); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusAccepted, map[string]string{"status": "demote requested"})
}

// resolveServer reads the {id} path parameter, fetches the server, and
// writes a 400/404 if anything is wrong. Returns ok=false if the handler
// must NOT continue (a response was already written).
func (d Deps) resolveServer(w http.ResponseWriter, r *http.Request) (*store.Server, bool) {
	id, ok := pathInt64(r, "id")
	if !ok {
		writeJSONError(w, http.StatusBadRequest, "invalid id")
		return nil, false
	}
	srv, err := d.Store.GetServer(id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeJSONError(w, http.StatusNotFound, "server not found")
			return nil, false
		}
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return nil, false
	}
	return srv, true
}
