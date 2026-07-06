package api

import (
	"context"
	"net/http"
	"strconv"
)

// handleServerTypes returns the Hetzner server types for the UI's
// server-type picker. It proxies to the first project's Hetzner API.
//
// Multi-tenant note: today the endpoint always reads from the first
// project in the store. A second project added later will not get its
// own catalog — the dropdown will reflect project #0's tokens. This is
// a known single-tenant limitation; the catalog is small (one Hetzner
// account is the typical deployment) so we accept it. If a future
// multi-tenant mode is added, the picker should switch to per-project
// token and a "current project" header.
func (d Deps) handleServerTypes(w http.ResponseWriter, r *http.Request) {
	types, err := d.serverTypes(r.Context())
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
// The error type is intentionally a tiny struct so the HTTP handler
// can pick the right status code without parsing strings; the metrics
// handler ignores it and falls back to the fixed map.
func (d Deps) serverTypes(ctx context.Context) ([]ServerTypeResponse, error) {
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
		// hcloud.Price.Gross is a string (e.g. "3.290000"); parse to
		// float32 for the DTO. If the string is unparseable, leave the
		// price at 0 rather than 500'ing the whole list.
		var priceMonthly float32
		if len(t.Pricings) > 0 {
			if v, perr := strconv.ParseFloat(t.Pricings[0].Monthly.Gross, 32); perr == nil {
				priceMonthly = float32(v)
			}
		}
		out = append(out, ServerTypeResponse{
			Name:            t.Name,
			Description:     t.Description,
			Cores:           t.Cores,
			MemoryGB:        t.Memory,
			DiskGB:          float32(t.Disk),
			Available:       len(t.Pricings) > 0,
			PriceMonthlyEUR: priceMonthly,
		})
	}
	return out, nil
}

// serverTypesError carries the HTTP status the original handler would
// have written, so the HTTP wrapper in handleServerTypes can reproduce
// the legacy status code without string-matching.
type serverTypesError struct {
	status int
	msg    string
}

func (e *serverTypesError) Error() string { return e.msg }