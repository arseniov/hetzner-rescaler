package api

import (
	"context"
	"net/http"
	"strconv"

	"github.com/jonamat/hetzner-rescaler/internal/hetzner"
)

// handleServerTypes returns Hetzner server types for the UI's
// server-type picker. It proxies to the first project's Hetzner API,
// filtered to the requested location's per-type availability.
//
// ?location=X is REQUIRED (400 otherwise) because availability and
// price are both per-location — returning the catalog without a
// location would be misleading.
//
// Multi-tenant note: today the endpoint always reads from the first
// project in the store. A second project added later will not get its
// own catalog — the dropdown will reflect project #0's tokens. This is
// a known single-tenant limitation; the catalog is small (one Hetzner
// account is the typical deployment) so we accept it. If a future
// multi-tenant mode is added, the picker should switch to per-project
// token and a "current project" header.
func (d Deps) handleServerTypes(w http.ResponseWriter, r *http.Request) {
	loc := r.URL.Query().Get("location")
	if loc == "" {
		writeJSONError(w, http.StatusBadRequest, "location query param required")
		return
	}
	types, err := d.serverTypes(r.Context(), loc)
	if err != nil {
		// serverTypes returns the raw error; map it to the right HTTP
		// status here (502 for upstream Hetzner failures, 500 for our
		// own). The pure function stays free of HTTP concerns so the
		// metrics handler can reuse it without wrapping in
		// httptest.NewRecorder.
		if errHe, ok := err.(*serverTypesError); ok {
			writeJSONError(w, errHe.status, errHe.msg)
			return
		}
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, types)
}

// serverTypes is the pure (no-HTTP) variant of handleServerTypes. It
// returns the same DTO slice the HTTP handler would have written, or
// an error. Used by both the HTTP handler and pricingMap so the
// metrics handler reads live prices from the same code path as the
// web's server-type picker.
//
// The location argument scopes both the "available" flag and the
// monthly EUR price to that location's entry on the underlying
// hetzner.ServerType. Types without any entry for the requested
// location come back Available=false, PriceMonthlyEUR=0.
//
// The error type is intentionally a tiny struct so the HTTP handler
// can pick the right status code without parsing strings; the metrics
// handler ignores it and falls back to the fixed map.
func (d Deps) serverTypes(ctx context.Context, location string) ([]ServerTypeResponse, error) {
	projects, err := d.Store.ListProjects()
	if err != nil {
		return nil, &serverTypesError{status: http.StatusInternalServerError, msg: err.Error()}
	}
	if len(projects) == 0 {
		return []ServerTypeResponse{}, nil
	}
	api, err := d.APIFor(projects[0].ID)
	if err != nil {
		return nil, &serverTypesError{status: http.StatusBadGateway, msg: "hetzner client: " + err.Error()}
	}
	types, err := api.ListServerTypes(ctx)
	if err != nil {
		return nil, &serverTypesError{status: http.StatusBadGateway, msg: "hetzner list types: " + err.Error()}
	}
	out := make([]ServerTypeResponse, 0, len(types))
	for _, t := range types {
		if t == nil {
			continue
		}
		out = append(out, ServerTypeResponse{
			Name:            t.Name,
			Description:     t.Description,
			Cores:           t.Cores,
			MemoryGB:        t.Memory,
			DiskGB:          float32(t.Disk),
			Available:       availableIn(t, location),
			PriceMonthlyEUR: priceIn(t, location),
		})
	}
	return out, nil
}

// availableIn returns true iff the given type has a Locations entry
// whose Location.Name matches and whose Available flag is set. A
// missing entry is treated as unavailable (false), since the type
// simply is not offered at the requested location.
func availableIn(t *hetzner.ServerType, location string) bool {
	for _, l := range t.Locations {
		if l.Location != nil && l.Location.Name == location {
			return l.Available
		}
	}
	return false
}

// priceIn parses the gross-eur monthly price from the Pricings entry
// whose Location.Name matches. Missing entry or unparseable Gross
// string yields 0.0 rather than an error — a $0 type is preferable
// to a 500 on the whole list.
func priceIn(t *hetzner.ServerType, location string) float32 {
	for _, p := range t.Pricings {
		if p.Location != nil && p.Location.Name == location {
			if v, perr := strconv.ParseFloat(p.Monthly.Gross, 32); perr == nil {
				return float32(v)
			}
		}
	}
	return 0
}

// serverTypesError carries the HTTP status the original handler would
// have written, so the HTTP wrapper in handleServerTypes can reproduce
// the legacy status code without string-matching.
type serverTypesError struct {
	status int
	msg    string
}

func (e *serverTypesError) Error() string { return e.msg }
