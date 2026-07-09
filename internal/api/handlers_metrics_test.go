package api

import (
	"context"
	"encoding/json"
	"math"
	"net/http"
	"testing"
	"time"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/jonamat/hetzner-rescaler/internal/hetzner"
	"github.com/jonamat/hetzner-rescaler/internal/store"
)

// TestMetricsAggregatesByDay verifies the plan's primary assertions:
//   - Rescales24hOk counts OK events in the last 24h
//   - LastRescaleError is populated when a recent failed event has an error
//   - RescaleCountsByDay has at least one row
//   - ActiveServerCount reflects seeded servers
func TestMetricsAggregatesByDay(t *testing.T) {
	deps, _ := newTestDeps(t)
	h := NewRouter(deps)

	// Two servers so ActiveServerCount is meaningful.
	_, sid1 := seedServer(t, deps, "p1", "web-1")
	_, sid2 := seedServer(t, deps, "p2", "web-2")

	now := time.Now().UTC()
	within24h := now.Add(-2 * time.Hour)
	day1 := now.Add(-25 * time.Hour) // day-1 (25h ago)
	day2 := now.Add(-49 * time.Hour) // day-2 (49h ago)

	// 3 OK + 1 fail on day-1 → contributes 3 to "Rescales24hOk" (since 25h ago is
	// already past the 24h cutoff, none of these count for 24h). We want
	// Rescales24hOk=3 so add 3 OK events inside 24h as well.
	events := []store.Event{
		{ServerID: sid1, Kind: "rescale_up", StartedAt: within24h, FinishedAt: within24h, OK: true, TriggeredBy: "scheduler"},
		{ServerID: sid1, Kind: "rescale_up", StartedAt: within24h, FinishedAt: within24h, OK: true, TriggeredBy: "scheduler"},
		{ServerID: sid2, Kind: "rescale_up", StartedAt: within24h, FinishedAt: within24h, OK: true, TriggeredBy: "scheduler"},
		{ServerID: sid1, Kind: "rescale_failed", StartedAt: day1, FinishedAt: day1, OK: false, Error: "boom", TriggeredBy: "scheduler"},
		{ServerID: sid1, Kind: "rescale_up", StartedAt: day1, FinishedAt: day1, OK: true, TriggeredBy: "scheduler"},
		{ServerID: sid1, Kind: "rescale_up", StartedAt: day1, FinishedAt: day1, OK: true, TriggeredBy: "scheduler"},
		{ServerID: sid2, Kind: "rescale_up", StartedAt: day1, FinishedAt: day1, OK: true, TriggeredBy: "scheduler"},
		{ServerID: sid1, Kind: "rescale_up", StartedAt: day2, FinishedAt: day2, OK: true, TriggeredBy: "scheduler"},
		{ServerID: sid1, Kind: "rescale_up", StartedAt: day2, FinishedAt: day2, OK: true, TriggeredBy: "scheduler"},
	}
	for _, e := range events {
		if _, err := deps.Store.AppendEvent(e); err != nil {
			t.Fatalf("AppendEvent: %v", err)
		}
	}

	req := authedRequest(t, http.MethodGet, "/api/metrics?range=7d", nil)
	rr := recorder(t, h, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("want 200, got %d (body=%q)", rr.Code, rr.Body.String())
	}

	var body MetricsResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode: %v", err)
	}

	if body.Range != "7d" {
		t.Errorf("range: got %q, want %q", body.Range, "7d")
	}
	if body.Kpis.Rescales24hOk != 3 {
		t.Errorf("Rescales24hOk: got %d, want 3", body.Kpis.Rescales24hOk)
	}
	if body.Kpis.LastRescaleError == nil {
		t.Fatalf("expected lastRescaleError to be set")
	}
	if body.Kpis.LastRescaleError.Error != "boom" {
		t.Errorf("lastRescaleError.Error: got %q, want %q", body.Kpis.LastRescaleError.Error, "boom")
	}
	if body.Kpis.ActiveServerCount != 2 {
		t.Errorf("ActiveServerCount: got %d, want 2", body.Kpis.ActiveServerCount)
	}
	if body.Kpis.ProjectsWithTokenCount != 2 {
		t.Errorf("ProjectsWithTokenCount: got %d, want 2", body.Kpis.ProjectsWithTokenCount)
	}
	if len(body.RescaleCountsByDay) == 0 {
		t.Errorf("expected non-empty RescaleCountsByDay")
	}
}

// TestMetricsRespectsRangeParam verifies that the response's "range"
// field reflects the ?range= query parameter.
func TestMetricsRespectsRangeParam(t *testing.T) {
	deps, _ := newTestDeps(t)
	h := NewRouter(deps)
	_, sid := seedServer(t, deps, "p1", "web-1")
	// Seed one event inside 24h so the 1d window is non-empty.
	now := time.Now().UTC()
	if _, err := deps.Store.AppendEvent(store.Event{
		ServerID: sid, Kind: "rescale_up",
		StartedAt: now.Add(-1 * time.Hour), FinishedAt: now.Add(-1 * time.Hour),
		OK: true, TriggeredBy: "scheduler",
	}); err != nil {
		t.Fatalf("AppendEvent: %v", err)
	}

	for _, q := range []string{"1d", "7d", "30d"} {
		t.Run(q, func(t *testing.T) {
			req := authedRequest(t, http.MethodGet, "/api/metrics?range="+q, nil)
			rr := recorder(t, h, req)
			if rr.Code != http.StatusOK {
				t.Fatalf("want 200, got %d", rr.Code)
			}
			var body MetricsResponse
			if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
				t.Fatalf("decode: %v", err)
			}
			if body.Range != q {
				t.Errorf("range: got %q, want %q", body.Range, q)
			}
		})
	}
}

// TestMetricsNoEventsReturnsEmpty verifies that with no events the
// endpoint returns valid arrays. RescaleCountsByDay is still populated
// with zero-rows for each day in the range (gaps-as-zeros); every
// bucket should have Total=0 and OK=Failed=0.
func TestMetricsNoEventsReturnsEmpty(t *testing.T) {
	deps, _ := newTestDeps(t)
	h := NewRouter(deps)
	_, _ = seedServer(t, deps, "p1", "web-1")

	req := authedRequest(t, http.MethodGet, "/api/metrics", nil)
	rr := recorder(t, h, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", rr.Code)
	}
	var body MetricsResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.Kpis.LastRescaleError != nil {
		t.Errorf("expected nil LastRescaleError, got %+v", body.Kpis.LastRescaleError)
	}
	if len(body.RescaleCountsByDay) == 0 {
		t.Errorf("expected day buckets for empty range, got 0")
	}
	for _, b := range body.RescaleCountsByDay {
		if b.Total != 0 || b.OK != 0 || b.Failed != 0 {
			t.Errorf("expected zero bucket %q, got %+v", b.Date, b)
		}
	}
	if len(body.HoursAtType) != 1 {
		t.Errorf("expected 1 HoursAtType row (one server seeded), got %d", len(body.HoursAtType))
	}
	if len(body.SuccessRateByServer) != 0 {
		t.Errorf("expected empty SuccessRateByServer, got %d rows", len(body.SuccessRateByServer))
	}
}

// TestMetricsRequiresAuth verifies the route is protected by the
// internal-token middleware.
func TestMetricsRequiresAuth(t *testing.T) {
	h := NewRouter(Deps{InternalToken: "secret"})
	req, err := http.NewRequest(http.MethodGet, "/api/metrics", nil)
	if err != nil {
		t.Fatalf("NewRequest: %v", err)
	}
	rr := recorder(t, h, req)
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("want 401, got %d", rr.Code)
	}
}

// TestPricingMap_PopulatesFromLiveTypes verifies that pricingMap pulls
// the per-location EUR price from the live Hetzner stub when the
// requested location has a matching pricing entry.
func TestPricingMap_PopulatesFromLiveTypes(t *testing.T) {
	deps, _ := newTestDeps(t)
	stub := &fakeHetzner{
		types: []*hetzner.ServerType{
			{
				Name: "cx11",
				Locations: []hcloud.ServerTypeLocation{{Location: &hcloud.Location{Name: "fsn1"}, Available: true}},
				Pricings: []hcloud.ServerTypeLocationPricing{
					{Location: &hcloud.Location{Name: "fsn1"}, Monthly: hcloud.Price{Currency: "EUR", Gross: "3.89"}},
				},
			},
		},
	}
	deps.APIFor = func(projectID int64) (hetzner.API, error) { return stub, nil }
	if _, err := deps.Store.CreateProject("p1", []byte("tok"), []byte("nonce12byts")); err != nil {
		t.Fatalf("CreateProject: %v", err)
	}

	pm := deps.pricingMap(context.Background(), "fsn1")
	// Assert the value matches the stub's live value, not just "> 0".
	// fixedPricingMap() already has cx11 at 3.29, so a > 0 assertion would
	// pass even if the live path were broken and the fallback map were
	// returned. 3.89 vs 3.29 proves the live lookup ran. We compare with
	// tolerance because float32→float64 round-tripping adds noise.
	got := pm["cx11"]
	if math.Abs(got-3.89) > 1e-3 {
		t.Fatalf("cx11 price = %v, want 3.89 (live stub value, not fixedPricingMap fallback 3.29)", got)
	}
}
