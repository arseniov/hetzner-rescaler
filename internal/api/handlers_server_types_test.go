package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/jonamat/hetzner-rescaler/internal/hetzner"
)

func TestServerTypes_ProxiesAndMarksUnavailable(t *testing.T) {
	deps, _ := newTestDeps(t)
	// hcloud's ServerType has no Available field. The handler derives
	// availability from Pricings: a sold-out type comes back with an
	// empty Pricings slice. cpx11 has a price (available); cpx31 does
	// not (unavailable).
	stub := &fakeHetzner{
		types: []*hetzner.ServerType{
			{
				Name: "cpx11", Description: "Intel/AMD shared", Cores: 2, Memory: 2.0, Disk: 40,
				Pricings: []hcloud.ServerTypeLocationPricing{{Location: &hcloud.Location{Name: "fsn1"}}},
			},
			{
				Name: "cpx31", Description: "Intel/AMD dedicated", Cores: 4, Memory: 8.0, Disk: 80,
				Pricings: nil,
			},
		},
	}
	deps.APIFor = func(projectID int64) (hetzner.API, error) { return stub, nil }
	// The handler picks the first project to know which Hetzner project to
	// query. Seed one (CreateProject needs encrypted token + nonce, same
	// pattern as seedServer).
	p, err := deps.Store.CreateProject("p1", []byte("tok"), []byte("nonce12byts"))
	if err != nil {
		t.Fatalf("CreateProject: %v", err)
	}
	_ = p
	h := NewRouter(deps)

	req := authedRequest(t, "GET", "/api/server-types", nil)
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
	if got[1].Available {
		t.Fatalf("cpx31 should be marked unavailable")
	}
}

func TestServerTypes_NoProjectsReturnsEmpty(t *testing.T) {
	deps, _ := newTestDeps(t)
	h := NewRouter(deps)
	req := authedRequest(t, "GET", "/api/server-types", nil)
	rr := recorder(t, h, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", rr.Code)
	}
	body := strings.TrimSpace(rr.Body.String())
	if body != "[]" {
		t.Fatalf("want [] for empty project list, got %q", body)
	}
}
