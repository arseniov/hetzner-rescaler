package api

import (
	"net/http"
	"strconv"
)

// handleServerTypes returns the Hetzner server types for the UI's
// server-type picker. It proxies to the first project's Hetzner API.
//
// The ServerTypeResponse DTO mirrors the hcloud SDK's ServerType:
// Name, Description, Cores, Memory (GB), Disk (GB), and the first
// pricing entry's monthly gross (EUR). hcloud's ServerType has no
// Available field, so availability is derived from Pricings: a
// sold-out type comes back with an empty Pricings slice; otherwise
// it has at least one location's monthly price. Hetzner returns one
// pricing entry per location; we use the first as a representative
// price for the picker.
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
	writeJSON(w, http.StatusOK, out)
}
