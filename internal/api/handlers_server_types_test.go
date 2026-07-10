package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/jonamat/hetzner-rescaler/internal/hetzner"
)

// TestServerTypes_ProxiesAndMarksUnavailable_PerLocation replaces the
// pre-location-gate version. It seeds two types with locations/pricings
// arrays for fsn1 + nbg1 and asserts that the per-location availability
// flag and price are derived correctly from the requested location.
func TestServerTypes_ProxiesAndMarksUnavailable_PerLocation(t *testing.T) {
	deps, _ := newTestDeps(t)
	stub := &fakeHetzner{
		types: []*hetzner.ServerType{
			{
				Name: "cpx11", Description: "Intel/AMD shared", Cores: 2, Memory: 2.0, Disk: 40,
				Locations: []hcloud.ServerTypeLocation{
					{Location: &hcloud.Location{Name: "fsn1"}, Available: true},
					{Location: &hcloud.Location{Name: "nbg1"}, Available: false},
				},
				Pricings: []hcloud.ServerTypeLocationPricing{
					{Location: &hcloud.Location{Name: "fsn1"}, Monthly: hcloud.Price{Currency: "EUR", Gross: "3.29"}},
					{Location: &hcloud.Location{Name: "nbg1"}, Monthly: hcloud.Price{Currency: "EUR", Gross: "3.50"}},
				},
			},
			{
				Name: "cpx31", Description: "dedicated", Cores: 4, Memory: 8.0, Disk: 80,
				Locations: []hcloud.ServerTypeLocation{
					{Location: &hcloud.Location{Name: "fsn1"}, Available: false},
					{Location: &hcloud.Location{Name: "nbg1"}, Available: false},
				},
				Pricings: nil,
			},
		},
	}
	deps.APIFor = func(projectID int64) (hetzner.API, error) { return stub, nil }
	if _, err := deps.Store.CreateProject("p1", []byte("tok"), []byte("nonce12byts")); err != nil {
		t.Fatalf("CreateProject: %v", err)
	}
	h := NewRouter(deps)

	req := authedRequest(t, "GET", "/api/server-types?location=fsn1", nil)
	rr := recorder(t, h, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("want 200, got %d (body=%q)", rr.Code, rr.Body.String())
	}
	var got []ServerTypeResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("want 2 types, got %d", len(got))
	}
	byName := map[string]ServerTypeResponse{}
	for _, t := range got {
		byName[t.Name] = t
	}
	if !byName["cpx11"].Available {
		t.Fatalf("cpx11 should be available in fsn1")
	}
	if byName["cpx11"].PriceMonthlyEUR <= 0 {
		t.Fatalf("cpx11 price should be the fsn1 price (3.29), got %v", byName["cpx11"].PriceMonthlyEUR)
	}
	if byName["cpx31"].Available {
		t.Fatalf("cpx31 should be unavailable in fsn1")
	}
	if byName["cpx31"].PriceMonthlyEUR != 0 {
		t.Fatalf("cpx31 price should be 0 when unavailable, got %v", byName["cpx31"].PriceMonthlyEUR)
	}
}

func TestServerTypes_MissingLocationReturns400(t *testing.T) {
	deps, _ := newTestDeps(t)
	h := NewRouter(deps)
	req := authedRequest(t, "GET", "/api/server-types", nil)
	rr := recorder(t, h, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d", rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "location") {
		t.Fatalf("body should mention location, got %q", rr.Body.String())
	}
}

func TestServerTypes_EmptyLocationReturns400(t *testing.T) {
	deps, _ := newTestDeps(t)
	h := NewRouter(deps)
	req := authedRequest(t, "GET", "/api/server-types?location=", nil)
	rr := recorder(t, h, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d", rr.Code)
	}
}

func TestServerTypes_TypeInOneLocationNotAnother(t *testing.T) {
	deps, _ := newTestDeps(t)
	stub := &fakeHetzner{
		types: []*hetzner.ServerType{
			{
				Name: "cpx11",
				Locations: []hcloud.ServerTypeLocation{
					{Location: &hcloud.Location{Name: "nbg1"}, Available: true},
				},
				Pricings: []hcloud.ServerTypeLocationPricing{
					{Location: &hcloud.Location{Name: "nbg1"}, Monthly: hcloud.Price{Currency: "EUR", Gross: "3.29"}},
				},
			},
		},
	}
	deps.APIFor = func(projectID int64) (hetzner.API, error) { return stub, nil }
	if _, err := deps.Store.CreateProject("p1", []byte("tok"), []byte("nonce12byts")); err != nil {
		t.Fatalf("CreateProject: %v", err)
	}
	h := NewRouter(deps)
	req := authedRequest(t, "GET", "/api/server-types?location=fsn1", nil)
	rr := recorder(t, h, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", rr.Code)
	}
	var got []ServerTypeResponse
	_ = json.Unmarshal(rr.Body.Bytes(), &got)
	if got[0].Available {
		t.Fatal("cpx11 should not be available in fsn1 (only available in nbg1)")
	}
	if got[0].PriceMonthlyEUR != 0 {
		t.Fatalf("cpx11 price should be 0 in fsn1 (no fsn1 pricing entry), got %v", got[0].PriceMonthlyEUR)
	}
}

func TestServerTypes_NoProjectsReturnsEmpty(t *testing.T) {
	deps, _ := newTestDeps(t)
	h := NewRouter(deps)
	req := authedRequest(t, "GET", "/api/server-types?location=fsn1", nil)
	rr := recorder(t, h, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", rr.Code)
	}
	if strings.TrimSpace(rr.Body.String()) != "[]" {
		t.Fatalf("want [], got %q", rr.Body.String())
	}
}
