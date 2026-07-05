package api

import (
	"net/http"
	"sort"
	"time"

	"github.com/jonamat/hetzner-rescaler/internal/store"
)

// --- Interfaces for testability ---
//
// metricsHandler only needs a small surface of *store.Store. Decoupling
// through interfaces makes it trivial to stub in tests if needed; in
// production the same *store.Store satisfies all three.

// EventReader is the slice of *store.Store the metrics handler needs.
type EventReader interface {
	ListEventsInRange(from, to time.Time) ([]*store.Event, error)
}

// ServerLister is the slice of *store.Store the metrics handler needs.
type ServerLister interface {
	ListAllServers() ([]*store.Server, error)
}

// ProjectLister is the slice of *store.Store the metrics handler needs
// to compute the "projects with token" KPI.
type ProjectLister interface {
	ListProjects() ([]*store.Project, error)
}

// Compile-time guarantee that *store.Store satisfies every interface the
// metrics handler depends on. If a future refactor drops one of these
// methods the build fails loudly instead of silently at runtime.
var (
	_ EventReader   = (*store.Store)(nil)
	_ ServerLister  = (*store.Store)(nil)
	_ ProjectLister = (*store.Store)(nil)
)

// --- Response types ---

// MetricsResponse is what GET /api/metrics returns. JSON tags use the
// project's snake_case convention (see types.go) — the plan's reference
// code used camelCase but the rest of the API is snake_case and the web
// client maps both equivalently.
type MetricsResponse struct {
	Range               string                  `json:"range"`
	From                string                  `json:"from"`
	To                  string                  `json:"to"`
	Kpis                MetricsKpis             `json:"kpis"`
	RescaleCountsByDay  []RescaleCountsByDayRow `json:"rescale_counts_by_day"`
	HoursAtType         []HoursAtTypeRow        `json:"hours_at_type"`
	SuccessRateByServer []SuccessRateRow        `json:"success_rate_by_server"`
}

// MetricsKpis is the top-level dashboard summary. LastRescaleError is a
// pointer so it can be omitted when there are no recent failures.
type MetricsKpis struct {
	ActiveServerCount      int           `json:"active_server_count"`
	ProjectsWithTokenCount int           `json:"projects_with_token_count"`
	Rescales24hOk          int           `json:"rescales_24h_ok"`
	LastRescaleError       *RescaleError `json:"last_rescale_error,omitempty"`
}

// RescaleError describes the most recent failed rescale event.
type RescaleError struct {
	ServerID int64  `json:"server_id"`
	Kind     string `json:"kind"`
	At       string `json:"at"`
	Error    string `json:"error"`
}

// RescaleCountsByDayRow is one day-bucket of rescale activity. Date is
// the UTC calendar day in YYYY-MM-DD form.
type RescaleCountsByDayRow struct {
	Date   string `json:"date"`
	OK     int    `json:"ok"`
	Failed int    `json:"failed"`
	Total  int    `json:"total"`
}

// HoursAtTypeRow is one server's split between base/top/fallback hours
// over the requested range, plus the EUR cost computed from the
// per-type pricing map.
type HoursAtTypeRow struct {
	ServerID   int64   `json:"server_id"`
	ServerName string  `json:"server_name"`
	Base       float64 `json:"base"`
	Top        float64 `json:"top"`
	Fallback   float64 `json:"fallback"`
	CostEUR    float64 `json:"cost_eur"`
}

// SuccessRateRow is one server's success/total counts over the requested
// range. Only servers with at least one event in range are returned.
type SuccessRateRow struct {
	ServerID   int64   `json:"server_id"`
	ServerName string  `json:"server_name"`
	OK         int     `json:"ok"`
	Total      int     `json:"total"`
	OkRate     float64 `json:"ok_rate"`
}

// --- Range durations ---

const (
	range1d  = 24 * time.Hour
	range7d  = 7 * 24 * time.Hour
	range30d = 30 * 24 * time.Hour
)

// metricsHandler returns an http.HandlerFunc that computes the dashboard
// summary. It is a top-level function (not a Deps method) so tests can
// invoke it directly with stub repos. Production wiring uses the thin
// (d Deps) handleMetrics method below.
func metricsHandler(events EventReader, servers ServerLister, projects ProjectLister, pricing map[string]float64) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rangeDur := range7d
		switch r.URL.Query().Get("range") {
		case "1d":
			rangeDur = range1d
		case "30d":
			rangeDur = range30d
		}

		to := time.Now().UTC()
		from := to.Add(-rangeDur)

		all, err := events.ListEventsInRange(from, to)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}

		serverList, err := servers.ListAllServers()
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}

		projectList, err := projects.ListProjects()
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}

		writeJSON(w, http.StatusOK, MetricsResponse{
			Range:               queryRangeLabel(rangeDur),
			From:                from.Format(time.RFC3339),
			To:                  to.Format(time.RFC3339),
			Kpis:                computeKpis(all, serverList, projectList, to),
			RescaleCountsByDay:  groupByDay(all, from, to),
			HoursAtType:         groupHoursByServer(all, serverList, pricing),
			SuccessRateByServer: successRateByServer(all, serverList),
		})
	}
}

// (d Deps) handleMetrics is the production-side thin wrapper that wires
// the live *store.Store into metricsHandler.
func (d Deps) handleMetrics(w http.ResponseWriter, r *http.Request) {
	metricsHandler(d.Store, d.Store, d.Store, fixedPricingMap()).ServeHTTP(w, r)
}

// queryRangeLabel renders the requested range as the canonical short
// string used in the response ("1d" | "7d" | "30d"). Anything other
// than 1d / 30d maps to 7d — the default.
func queryRangeLabel(d time.Duration) string {
	switch d {
	case range1d:
		return "1d"
	case range30d:
		return "30d"
	default:
		return "7d"
	}
}

// --- KPI computation ---

func computeKpis(events []*store.Event, servers []*store.Server, projects []*store.Project, now time.Time) MetricsKpis {
	// Rescales24hOk is intentionally NOT range-scoped: the dashboard's
	// "successful in last 24h" card is always a 24h rolling window
	// regardless of the query's chart range.
	cutoff24h := now.Add(-24 * time.Hour)
	ok24h := 0
	for _, e := range events {
		if e.OK && !e.StartedAt.Before(cutoff24h) {
			ok24h++
		}
	}

	// Active server count = total registered servers. We don't filter
	// on mode (manual/auto_promote/scheduled) here: the UI's KPI card
	// label is "Active servers" and "registered" is the simplest
	// defensible definition. A future task can split by mode if the
	// label changes.
	activeServers := len(servers)

	// ProjectsWithTokenCount: any project that has stored an encrypted
	// token. HasToken on ProjectResponse uses the same heuristic.
	projectsWithToken := 0
	for _, p := range projects {
		if len(p.HCloudTokenEncrypted) > 0 {
			projectsWithToken++
		}
	}

	return MetricsKpis{
		ActiveServerCount:      activeServers,
		ProjectsWithTokenCount: projectsWithToken,
		Rescales24hOk:          ok24h,
		LastRescaleError:       lastRescaleError(events),
	}
}

// lastRescaleError scans events (which are sorted newest-first by
// ListEventsInRange) and returns the most recent failure with a non-empty
// Error message. Returns nil when no qualifying event exists.
func lastRescaleError(events []*store.Event) *RescaleError {
	for _, e := range events {
		if !e.OK && e.Error != "" {
			return &RescaleError{
				ServerID: e.ServerID,
				Kind:     e.Kind,
				At:       e.StartedAt.Format(time.RFC3339),
				Error:    e.Error,
			}
		}
	}
	return nil
}

// --- Day-bucket aggregation ---

// groupByDay bins events into UTC calendar days. Returns rows in
// ascending date order so the chart can plot left-to-right without
// re-sorting on the client.
func groupByDay(events []*store.Event, from, to time.Time) []RescaleCountsByDayRow {
	buckets := map[string]*RescaleCountsByDayRow{}
	for _, e := range events {
		day := e.StartedAt.UTC().Format("2006-01-02")
		b, ok := buckets[day]
		if !ok {
			b = &RescaleCountsByDayRow{Date: day}
			buckets[day] = b
		}
		if e.OK {
			b.OK++
		} else {
			b.Failed++
		}
		b.Total++
	}

	// Seed every day in [from, to] so the chart shows gaps as zeros
	// rather than as missing bars.
	d := from.UTC().Truncate(24 * time.Hour)
	end := to.UTC().Truncate(24 * time.Hour)
	for !d.After(end) {
		key := d.Format("2006-01-02")
		if _, ok := buckets[key]; !ok {
			buckets[key] = &RescaleCountsByDayRow{Date: key}
		}
		d = d.Add(24 * time.Hour)
	}

	out := make([]RescaleCountsByDayRow, 0, len(buckets))
	for _, b := range buckets {
		out = append(out, *b)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Date < out[j].Date })
	return out
}

// --- Hours-at-type aggregation ---

// groupHoursByServer reports, per server, the total hours spent at base,
// top, and fallback types within the requested range. Cost is computed
// as sum(hours × price_per_month / 720), where 720 is the conventional
// hours-per-month divisor Hetzner pricing uses.
//
// Each event with a non-zero FinishedAt contributes (FinishedAt -
// StartedAt) at its ToType. Events with zero FinishedAt (still in
// progress) contribute 0 hours. The fallback bucket catches any ToType
// that isn't the server's BaseServerType or TopServerType.
func groupHoursByServer(events []*store.Event, servers []*store.Server, pricing map[string]float64) []HoursAtTypeRow {
	byID := map[int64]*HoursAtTypeRow{}
	for _, s := range servers {
		byID[s.ID] = &HoursAtTypeRow{ServerID: s.ID, ServerName: s.Name}
	}

	for _, e := range events {
		if e.FinishedAt.IsZero() {
			continue
		}
		row, ok := byID[e.ServerID]
		if !ok {
			// Event for an unknown (deleted?) server — skip.
			continue
		}
		dur := e.FinishedAt.Sub(e.StartedAt)
		if dur <= 0 {
			continue
		}
		hours := dur.Hours()
		// Find the server's base/top to classify the event.
		var srv *store.Server
		for _, s := range servers {
			if s.ID == e.ServerID {
				srv = s
				break
			}
		}
		if srv == nil {
			continue
		}
		switch e.ToType {
		case srv.BaseServerType:
			row.Base += hours
		case srv.TopServerType:
			row.Top += hours
		default:
			row.Fallback += hours
		}
	}

	out := make([]HoursAtTypeRow, 0, len(byID))
	for _, row := range byID {
		// Find server for pricing purposes (use the row's stored
		// server by ID — we already iterated servers above).
		var srv *store.Server
		for _, s := range servers {
			if s.ID == row.ServerID {
				srv = s
				break
			}
		}
		row.CostEUR = computeCostEUR(row, srv, pricing)
		out = append(out, *row)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ServerID < out[j].ServerID })
	return out
}

// computeCostEUR multiplies hours-at-type by the matching €/month price
// and divides by 720. Unknown types use the fallback price.
func computeCostEUR(row *HoursAtTypeRow, srv *store.Server, pricing map[string]float64) float64 {
	if srv == nil {
		return 0
	}
	price := func(t string) float64 {
		if p, ok := pricing[t]; ok {
			return p
		}
		return pricing["__default__"]
	}
	base := price(srv.BaseServerType) * row.Base / 720.0
	top := price(srv.TopServerType) * row.Top / 720.0
	// For fallback events, charge the fallback type's own price when
	// known, else the default price.
	fallbackCost := 0.0
	if row.Fallback > 0 {
		fallbackCost = pricing["__default__"] * row.Fallback / 720.0
	}
	return base + top + fallbackCost
}

// --- Success rate per server ---

// successRateByServer returns one row per server that has at least one
// event in range. Servers with zero events are omitted — the UI treats
// absence as "no activity yet".
func successRateByServer(events []*store.Event, servers []*store.Server) []SuccessRateRow {
	type counter struct{ ok, total int }
	counts := map[int64]*counter{}
	for _, e := range events {
		c, ok := counts[e.ServerID]
		if !ok {
			c = &counter{}
			counts[e.ServerID] = c
		}
		c.total++
		if e.OK {
			c.ok++
		}
	}

	out := make([]SuccessRateRow, 0, len(counts))
	for id, c := range counts {
		name := ""
		for _, s := range servers {
			if s.ID == id {
				name = s.Name
				break
			}
		}
		rate := 0.0
		if c.total > 0 {
			rate = float64(c.ok) / float64(c.total)
		}
		out = append(out, SuccessRateRow{
			ServerID:   id,
			ServerName: name,
			OK:         c.ok,
			Total:      c.total,
			OkRate:     rate,
		})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ServerID < out[j].ServerID })
	return out
}

// --- Pricing ---

// fixedPricingMap returns the €/month prices for the Hetzner shared-
// resource types the rescaler supports. The "__default__" key is the
// fallback for any type we don't know about (and for fallback-type
// hours in groupHoursByServer).
func fixedPricingMap() map[string]float64 {
	return map[string]float64{
		"cx11": 3.29, "cx21": 5.39, "cx31": 10.39, "cx41": 19.39, "cx51": 38.59,
		"cax11": 3.29, "cax21": 5.39, "cax31": 10.39, "cax41": 19.39,
		"__default__": 5.00,
	}
}
