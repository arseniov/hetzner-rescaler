package api

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/jonamat/hetzner-rescaler/internal/hetzner"
	"github.com/jonamat/hetzner-rescaler/internal/store"
)

func seedServer(t *testing.T, deps Deps, projectName, serverName string) (int64, int64) {
	t.Helper()
	p, err := deps.Store.CreateProject(projectName, []byte("tok"), []byte("nonce12byts"))
	if err != nil {
		t.Fatalf("CreateProject: %v", err)
	}
	s, err := deps.Store.CreateServer(p.ID, store.Server{
		HCloudServerID: 1, Name: serverName, Label: serverName,
		BaseServerType: "cpx11", TopServerType: "cpx31",
		FallbackChain:  []string{"cpx31", "cpx11"},
		Mode: "manual", Timezone: "UTC",
	})
	if err != nil {
		t.Fatalf("CreateServer: %v", err)
	}
	return p.ID, s.ID
}

func TestListServers_ReturnsArrayAcrossProjects(t *testing.T) {
	deps, _ := newTestDeps(t)
	h := NewRouter(deps)
	_, s1 := seedServer(t, deps, "p1", "web-1")
	_, s2 := seedServer(t, deps, "p2", "web-2")

	req := authedRequest(t, "GET", "/api/servers", nil)
	rr := recorder(t, h, req)

	var got []ServerResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("want 2 servers, got %d", len(got))
	}
	ids := map[int64]bool{s1: false, s2: false}
	for _, s := range got {
		ids[s.ID] = true
	}
	for id, seen := range ids {
		if !seen {
			t.Fatalf("missing server id %d", id)
		}
	}
}

func TestGetServer_Returns404WhenMissing(t *testing.T) {
	deps, _ := newTestDeps(t)
	h := NewRouter(deps)
	req := authedRequest(t, "GET", "/api/servers/9999", nil)
	rr := recorder(t, h, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("want 404, got %d", rr.Code)
	}
}

func TestGetServer_ReturnsServer(t *testing.T) {
	deps, _ := newTestDeps(t)
	h := NewRouter(deps)
	_, sid := seedServer(t, deps, "p1", "web-1")
	req := authedRequest(t, "GET", "/api/servers/"+itoa(sid), nil)
	rr := recorder(t, h, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("want 200, got %d (body=%q)", rr.Code, rr.Body.String())
	}
	var got ServerResponse
	_ = json.Unmarshal(rr.Body.Bytes(), &got)
	if got.ID != sid || got.Name != "web-1" {
		t.Fatalf("unexpected: %+v", got)
	}
}

func TestCreateServer_Validates(t *testing.T) {
	deps, _ := newTestDeps(t)
	h := NewRouter(deps)
	pid, _ := seedServer(t, deps, "p1", "web-1")

	cases := []CreateServerRequest{
		{ProjectID: pid, HCloudServerID: 2, Name: "x"}, // missing fields
		{ProjectID: 999, HCloudServerID: 2, Name: "x", BaseServerType: "cpx11",
			TopServerType: "cpx31", FallbackChain: []string{"cpx31"}, Mode: "manual", Timezone: "UTC"},
	}
	for i, body := range cases {
		req := authedRequest(t, "POST", "/api/servers", body)
		rr := recorder(t, h, req)
		if rr.Code == http.StatusCreated {
			t.Fatalf("case %d: expected non-201, got 201", i)
		}
	}
}

func TestCreateServer_Succeeds(t *testing.T) {
	deps, _ := newTestDeps(t)
	h := NewRouter(deps)
	pid, _ := seedServer(t, deps, "p1", "web-1")
	body := CreateServerRequest{
		ProjectID: pid, HCloudServerID: 99, Name: "web-99", Label: "secondary",
		BaseServerType: "cpx11", TopServerType: "cpx31",
		FallbackChain: []string{"cpx31", "cpx11"}, Mode: "scheduled", Timezone: "UTC",
	}
	req := authedRequest(t, "POST", "/api/servers", body)
	rr := recorder(t, h, req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("want 201, got %d (body=%q)", rr.Code, rr.Body.String())
	}
	var got ServerResponse
	_ = json.Unmarshal(rr.Body.Bytes(), &got)
	if got.Name != "web-99" || got.Mode != "scheduled" {
		t.Fatalf("unexpected: %+v", got)
	}
}

func TestUpdateServer_AppliesChanges(t *testing.T) {
	deps, _ := newTestDeps(t)
	h := NewRouter(deps)
	_, sid := seedServer(t, deps, "p1", "web-1")
	body := UpdateServerRequest{
		Name: "web-1-renamed", Label: "primary",
		BaseServerType: "cpx11", TopServerType: "cpx31",
		FallbackChain: []string{"cpx31", "cpx11"}, Mode: "auto_promote", Timezone: "UTC",
	}
	req := authedRequest(t, "PUT", "/api/servers/"+itoa(sid), body)
	rr := recorder(t, h, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("want 200, got %d (body=%q)", rr.Code, rr.Body.String())
	}
	var got ServerResponse
	_ = json.Unmarshal(rr.Body.Bytes(), &got)
	if got.Name != "web-1-renamed" || got.Mode != "auto_promote" {
		t.Fatalf("update did not apply: %+v", got)
	}
}

func TestDeleteServer_RemovesRow(t *testing.T) {
	deps, _ := newTestDeps(t)
	h := NewRouter(deps)
	_, sid := seedServer(t, deps, "p1", "web-1")
	req := authedRequest(t, "DELETE", "/api/servers/"+itoa(sid), nil)
	rr := recorder(t, h, req)
	if rr.Code != http.StatusNoContent {
		t.Fatalf("want 204, got %d", rr.Code)
	}
	req = authedRequest(t, "GET", "/api/servers/"+itoa(sid), nil)
	rr = recorder(t, h, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("want 404 after delete, got %d", rr.Code)
	}
}

// newLiveServer is a small constructor for hcloud.Server values used
// only in tests for the live-state plumbing. Keeping it in the test
// file avoids polluting the main package with a constructor that
// has no production call site.
func newLiveServer(id int, status hcloud.ServerStatus, typeName string) *hetzner.Server {
	return &hetzner.Server{
		ID:         id,
		Name:       "live-" + itoa(int64(id)),
		Status:     status,
		ServerType: &hetzner.ServerType{Name: typeName},
	}
}

func TestListServers_PopulatesLiveState(t *testing.T) {
	deps, _ := newTestDeps(t)
	_, _ = seedServer(t, deps, "p1", "web-1")
	// The seeded server has HCloudServerID=1. Make our fake return a
	// running server with type "cx42" so we can assert the live fields
	// flow through to the JSON response.
	stub := &fakeHetzner{
		servers: []*hetzner.Server{
			newLiveServer(1, hcloud.ServerStatusRunning, "cx42"),
		},
	}
	deps.APIFor = func(projectID int64) (hetzner.API, error) { return stub, nil }

	h := NewRouter(deps)
	req := authedRequest(t, "GET", "/api/servers", nil)
	rr := recorder(t, h, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("want 200, got %d (body=%q)", rr.Code, rr.Body.String())
	}
	var got []ServerResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("want 1 server, got %d", len(got))
	}
	if got[0].Status != "running" {
		t.Fatalf("status not populated: got %q", got[0].Status)
	}
	if got[0].CurrentType != "cx42" {
		t.Fatalf("current_type not populated: got %q", got[0].CurrentType)
	}
	// created_at / updated_at must round-trip — they're sourced from
	// the store, not the live API.
	if got[0].CreatedAt.IsZero() {
		t.Fatalf("created_at is zero")
	}
}

func TestListServers_OmitsLiveStateWhenHetznerFails(t *testing.T) {
	deps, _ := newTestDeps(t)
	_, _ = seedServer(t, deps, "p1", "web-1")
	// Stub that has no servers: liveServerState returns zero-value
	// (GetServer returns (nil, nil) → hs == nil branch fires).
	stub := &fakeHetzner{}
	deps.APIFor = func(projectID int64) (hetzner.API, error) { return stub, nil }

	h := NewRouter(deps)
	req := authedRequest(t, "GET", "/api/servers", nil)
	rr := recorder(t, h, req)
	if rr.Code != http.StatusOK {
		// A Hetzner failure must NOT turn into a 500. The endpoint's
		// contract is "store data + best-effort live data"; degradation
		// is the only acceptable behavior here.
		t.Fatalf("want 200 with soft-fail, got %d (body=%q)", rr.Code, rr.Body.String())
	}
	// Decode into a map so we can prove the live fields are absent
	// rather than just empty strings.
	var raw []map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &raw); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(raw) != 1 {
		t.Fatalf("want 1 server, got %d", len(raw))
	}
	if _, ok := raw[0]["status"]; ok {
		t.Fatalf("status should be omitted on soft-fail; got %v", raw[0]["status"])
	}
	if _, ok := raw[0]["current_type"]; ok {
		t.Fatalf("current_type should be omitted on soft-fail; got %v", raw[0]["current_type"])
	}
}

func TestGetServer_PopulatesLiveState(t *testing.T) {
	deps, _ := newTestDeps(t)
	_, sid := seedServer(t, deps, "p1", "web-1")
	stub := &fakeHetzner{
		servers: []*hetzner.Server{
			newLiveServer(1, hcloud.ServerStatusInitializing, "cpx31"),
		},
	}
	deps.APIFor = func(projectID int64) (hetzner.API, error) { return stub, nil }

	h := NewRouter(deps)
	req := authedRequest(t, "GET", "/api/servers/"+itoa(sid), nil)
	rr := recorder(t, h, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", rr.Code)
	}
	var got ServerResponse
	_ = json.Unmarshal(rr.Body.Bytes(), &got)
	if got.Status != "initializing" {
		t.Fatalf("status: got %q want initializing", got.Status)
	}
	if got.CurrentType != "cpx31" {
		t.Fatalf("current_type: got %q want cpx31", got.CurrentType)
	}
}

// TestLiveServerState_SoftFailOnNilDepsAPIFor guards the contract that
// a misconfigured Deps (no APIFor) cannot panic the handler.
func TestLiveServerState_SoftFailOnNilDepsAPIFor(t *testing.T) {
	d := Deps{APIFor: nil}
	srv := &store.Server{ProjectID: 1, HCloudServerID: 1}
	got := d.liveServerState(context.Background(), srv)
	if got != (LiveServerState{}) {
		t.Fatalf("expected zero LiveServerState, got %+v", got)
	}
}

func TestGetServer_IncludesPendingEvent(t *testing.T) {
	deps, _ := newTestDeps(t)
	_, sid := seedServer(t, deps, "p1", "web-1")

	// Insert a pending event.
	_, err := deps.Store.AppendEvent(store.Event{
		ServerID: sid, Kind: "rescale_pending",
		StartedAt: time.Now().UTC(), TriggeredBy: "test",
	})
	if err != nil {
		t.Fatalf("AppendEvent: %v", err)
	}

	h := NewRouter(deps)
	req := authedRequest(t, "GET", "/api/servers/"+itoa(sid), nil)
	rr := recorder(t, h, req)

	var got ServerResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.PendingEvent == nil {
		t.Fatal("pending_event missing from response")
	}
	if got.PendingEvent.Kind != "rescale_pending" {
		t.Fatalf("pending_event.kind = %q, want rescale_pending", got.PendingEvent.Kind)
	}
}

func TestGetServer_OmitsPendingEventWhenNone(t *testing.T) {
	deps, _ := newTestDeps(t)
	_, sid := seedServer(t, deps, "p1", "web-1")

	h := NewRouter(deps)
	req := authedRequest(t, "GET", "/api/servers/"+itoa(sid), nil)
	rr := recorder(t, h, req)

	var got map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if _, ok := got["pending_event"]; ok {
		t.Fatalf("pending_event should be absent, got %v", got["pending_event"])
	}
}