package api

import (
	"encoding/json"
	"net/http"
	"testing"

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