package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jonamat/hetzner-rescaler/internal/rescaler"
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
	if d.Manager == nil {
		writeJSONError(w, http.StatusServiceUnavailable, "rescaler manager not configured")
		return
	}
	id, err := d.Manager.Submit(r.Context(), srv, target, "api")
	if err != nil {
		if errors.Is(err, rescaler.ErrAlreadyInProgress) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			_ = json.NewEncoder(w).Encode(map[string]any{
				"error":            "rescale already in progress",
				"pending_event_id": d.activePendingID(srv.ID),
			})
			return
		}
		writeJSONError(w, http.StatusInternalServerError, "rescale submit failed: "+err.Error())
		return
	}
	writeJSON(w, http.StatusAccepted, map[string]any{
		"status":           "rescale initiated",
		"target":           target,
		"pending_event_id": id,
	})
}

// activePendingID returns the ID of the still-pending rescale_pending
// event for srv, or 0 if there is none (or if the lookup fails).
// Convenience wrapper used by handlers that want to surface the pending
// ID alongside other response fields.
func (d Deps) activePendingID(srvID int64) int64 {
	pending, err := d.Store.ActivePendingEventForServer(srvID)
	if err != nil || pending == nil {
		return 0
	}
	return pending.ID
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
	// UpdateServer takes srv by value, so mutating the local pointer is safe.
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
	// UpdateServer takes srv by value, so mutating the local pointer is safe.
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
