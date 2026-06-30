package api

import "net/http"

// handleServerTypes returns the Hetzner server types for the UI's
// server-type picker. It proxies to the first project's Hetzner API.
//
// The plan's "ServerTypeResponse" DTO has more fields (Description, Cores,
// MemoryGB, DiskGB, PriceMonthlyEUR) than this minimal handler maps —
// see Task 9 in the plan. Per the plan scope we only map Name and
// Available. Future tasks can expand the mapping if the SPA needs more.
func (d Deps) handleServerTypes(w http.ResponseWriter, r *http.Request) {
	projects, err := d.Store.ListProjects()
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if len(projects) == 0 {
		writeJSON(w, http.StatusOK, []ServerTypeResponse{})
		return
	}
	api, err := d.APIFor(projects[0].ID)
	if err != nil {
		writeJSONError(w, http.StatusBadGateway, "hetzner client: "+err.Error())
		return
	}
	types, err := api.ListServerTypes(r.Context())
	if err != nil {
		writeJSONError(w, http.StatusBadGateway, "hetzner list types: "+err.Error())
		return
	}
	// hcloud's ServerType has no Available field. We infer availability
	// from Pricings: a sold-out type comes back with an empty Pricings
	// slice; otherwise it has at least one location's monthly price.
	out := make([]ServerTypeResponse, 0, len(types))
	for _, t := range types {
		if t == nil {
			continue
		}
		out = append(out, ServerTypeResponse{
			Name:      t.Name,
			Available: len(t.Pricings) > 0,
		})
	}
	writeJSON(w, http.StatusOK, out)
}
