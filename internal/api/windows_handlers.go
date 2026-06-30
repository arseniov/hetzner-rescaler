package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"time"

	"github.com/jonamat/hetzner-rescaler/internal/store"
)

var hhmmRe = regexp.MustCompile(`^([01][0-9]|2[0-3]):[0-5][0-9]$`)

func (d Deps) handleListWindows(w http.ResponseWriter, r *http.Request) {
	sid, ok := pathInt64(r, "id")
	if !ok {
		writeJSONError(w, http.StatusBadRequest, "invalid id")
		return
	}
	wins, err := d.Store.ListWindows(sid)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	out := make([]WindowResponse, 0, len(wins))
	for _, win := range wins {
		out = append(out, windowToResponse(win))
	}
	writeJSON(w, http.StatusOK, out)
}

func (d Deps) handleCreateWindow(w http.ResponseWriter, r *http.Request) {
	sid, ok := pathInt64(r, "id")
	if !ok {
		writeJSONError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var req CreateWindowRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if err := validateWindowRequest(req); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	win, err := d.Store.CreateWindow(sid, store.Window{
		Label:      req.Label,
		DaysOfWeek: req.DaysOfWeek,
		StartTime:  req.StartTime,
		StopTime:   req.StopTime,
		TargetType: req.TargetType,
		Enabled:    req.Enabled,
	})
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, windowToResponse(win))
}

func (d Deps) handleUpdateWindow(w http.ResponseWriter, r *http.Request) {
	wid, ok := pathInt64(r, "wid")
	if !ok {
		writeJSONError(w, http.StatusBadRequest, "invalid wid")
		return
	}
	var req UpdateWindowRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if err := validateWindowRequest(req); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	win, err := d.updateWindow(wid, req)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeJSONError(w, http.StatusNotFound, "window not found")
			return
		}
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, windowToResponse(win))
}

func (d Deps) handleDeleteWindow(w http.ResponseWriter, r *http.Request) {
	wid, ok := pathInt64(r, "wid")
	if !ok {
		writeJSONError(w, http.StatusBadRequest, "invalid wid")
		return
	}
	if err := d.Store.DeleteWindow(wid); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// validateWindowRequest checks the request shape. Returns a human-readable
// error or nil.
func validateWindowRequest(req CreateWindowRequest) error {
	if req.Label == "" {
		return errors.New("label is required")
	}
	if req.TargetType == "" {
		return errors.New("target_type is required")
	}
	if !hhmmRe.MatchString(req.StartTime) || !hhmmRe.MatchString(req.StopTime) {
		return errors.New("start_time and stop_time must be HH:MM (24h)")
	}
	startMins := hhmmToMinutes(req.StartTime)
	stopMins := hhmmToMinutes(req.StopTime)
	if startMins >= stopMins {
		return errors.New("start_time must be < stop_time")
	}
	if req.DaysOfWeek < 0 || req.DaysOfWeek > 0b1111111 {
		return errors.New("days_of_week must be a 7-bit bitmask")
	}
	return nil
}

func hhmmToMinutes(s string) int {
	t, _ := time.Parse("15:04", s)
	return t.Hour()*60 + t.Minute()
}

// windowToResponse converts a *store.Window to its API projection.
func windowToResponse(w *store.Window) WindowResponse {
	return WindowResponse{
		ID:         w.ID,
		ServerID:   w.ServerID,
		Label:      w.Label,
		DaysOfWeek: w.DaysOfWeek,
		StartTime:  w.StartTime,
		StopTime:   w.StopTime,
		TargetType: w.TargetType,
		Enabled:    w.Enabled,
	}
}

// updateWindow invokes the store's UpdateWindow method. Kept here so
// handleUpdateWindow stays thin and shared validation logic lives in one
// place.
func (d Deps) updateWindow(id int64, req UpdateWindowRequest) (*store.Window, error) {
	return d.Store.UpdateWindow(id, store.Window{
		Label:      req.Label,
		DaysOfWeek: req.DaysOfWeek,
		StartTime:  req.StartTime,
		StopTime:   req.StopTime,
		TargetType: req.TargetType,
		Enabled:    req.Enabled,
	})
}
