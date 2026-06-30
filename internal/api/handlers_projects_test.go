package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jonamat/hetzner-rescaler/internal/crypto"
	"github.com/jonamat/hetzner-rescaler/internal/hetzner"
	"github.com/jonamat/hetzner-rescaler/internal/store"
)

const testInternalToken = "test-internal-token"

// newTestDeps builds a Deps with a real (temp) SQLite store, an in-memory
// AES-GCM key, and a stub Hetzner API factory. Tests can override APIFor
// per-subtest.
func newTestDeps(t *testing.T) (Deps, *crypto.Keyring) {
	t.Helper()
	s, err := store.OpenTemp()
	if err != nil {
		t.Fatalf("store.OpenTemp: %v", err)
	}
	t.Cleanup(func() { s.Close() })

	k, err := crypto.NewKeyring()
	if err != nil {
		t.Fatalf("crypto.NewKeyring: %v", err)
	}
	return Deps{
		InternalToken: testInternalToken,
		Store:         s,
		Keyring:       k,
		APIFor: func(projectID int64) (hetzner.API, error) {
			return nil, nil // overridden in tests that need Hetzner
		},
	}, k
}

// authedRequest is defined in testhelpers_test.go.

func authedRequest(t *testing.T, method, path string, body any) *http.Request {
	t.Helper()
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			t.Fatalf("encode body: %v", err)
		}
	}
	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set("X-Internal-Token", testInternalToken)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return req
}

func TestListProjects_EmptyReturnsEmptyArray(t *testing.T) {
	deps, _ := newTestDeps(t)
	h := NewRouter(deps)

	req := authedRequest(t, http.MethodGet, "/api/projects", nil)
	rr := recorder(t, h, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("want 200, got %d (body=%q)", rr.Code, rr.Body.String())
	}
	body := strings.TrimSpace(rr.Body.String())
	if body != "[]" {
		t.Fatalf("want [] for empty list, got %q", body)
	}
}

func TestCreateProject_PersistsAndEncryptsToken(t *testing.T) {
	deps, _ := newTestDeps(t)
	h := NewRouter(deps)

	body := CreateProjectRequest{Name: "prod", HCloudToken: "hcloud-secret-123"}
	req := authedRequest(t, http.MethodPost, "/api/projects", body)
	rr := recorder(t, h, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("want 201, got %d (body=%q)", rr.Code, rr.Body.String())
	}
	var got ProjectResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.ID == 0 {
		t.Fatalf("want non-zero ID, got %d", got.ID)
	}
	if got.Name != "prod" {
		t.Fatalf("want name=prod, got %q", got.Name)
	}
	if !got.HasToken {
		t.Fatalf("want has_token=true")
	}

	// Re-list and confirm the token is NOT in the response.
	req = authedRequest(t, http.MethodGet, "/api/projects", nil)
	rr = recorder(t, h, req)
	if !bytes.Contains(rr.Body.Bytes(), []byte(`"has_token":true`)) {
		t.Fatalf("want has_token=true in list, got %s", rr.Body.String())
	}
	if bytes.Contains(rr.Body.Bytes(), []byte("hcloud-secret-123")) {
		t.Fatalf("token leaked in list response")
	}
}

func TestCreateProject_RejectsMissingFields(t *testing.T) {
	deps, _ := newTestDeps(t)
	h := NewRouter(deps)

	cases := []struct {
		name string
		body CreateProjectRequest
	}{
		{"empty name", CreateProjectRequest{Name: "", HCloudToken: "x"}},
		{"empty token", CreateProjectRequest{Name: "x", HCloudToken: ""}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := authedRequest(t, http.MethodPost, "/api/projects", tc.body)
			rr := recorder(t, h, req)
			if rr.Code != http.StatusBadRequest {
				t.Fatalf("want 400, got %d", rr.Code)
			}
		})
	}
}

func TestDeleteProject_RemovesRow(t *testing.T) {
	deps, _ := newTestDeps(t)
	h := NewRouter(deps)

	body := CreateProjectRequest{Name: "prod", HCloudToken: "tok"}
	req := authedRequest(t, http.MethodPost, "/api/projects", body)
	rr := recorder(t, h, req)
	var created ProjectResponse
	_ = json.Unmarshal(rr.Body.Bytes(), &created)

	req = authedRequest(t, http.MethodDelete, "/api/projects/"+itoa(created.ID), nil)
	rr = recorder(t, h, req)
	if rr.Code != http.StatusNoContent {
		t.Fatalf("want 204, got %d", rr.Code)
	}

	// Second delete returns 404.
	req = authedRequest(t, http.MethodDelete, "/api/projects/"+itoa(created.ID), nil)
	rr = recorder(t, h, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("want 404 on second delete, got %d", rr.Code)
	}
}

func TestRefreshProject_AddsNewServers(t *testing.T) {
	deps, _ := newTestDeps(t)

	// Seed a project. The handler will encrypt via Deps.Keyring when it
	// re-uses the encrypted form on refresh, but for the seed we just
	// need a project row that exists.
	created, err := deps.Store.CreateProject("prod", []byte("tok-enc"), []byte("nonce12byts"))
	if err != nil {
		t.Fatalf("CreateProject: %v", err)
	}

	// Stub Hetzner: returns two servers. Test the "add new" path only.
	stub := &fakeHetzner{
		servers: []*hetzner.Server{
			{ID: 1, Name: "web-1", ServerType: &hetzner.ServerType{Name: "cpx11"}},
			{ID: 2, Name: "web-2", ServerType: &hetzner.ServerType{Name: "cpx11"}},
		},
	}
	deps.APIFor = func(projectID int64) (hetzner.API, error) { return stub, nil }
	h := NewRouter(deps)

	req := authedRequest(t, http.MethodPost, "/api/projects/"+itoa(created.ID)+"/refresh", nil)
	rr := recorder(t, h, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("want 200, got %d (body=%q)", rr.Code, rr.Body.String())
	}
	var resp RefreshProjectResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(resp.Added) != 2 {
		t.Fatalf("want 2 added, got %d", len(resp.Added))
	}
}

func TestRefreshProject_SkipsAlreadyRegistered(t *testing.T) {
	deps, _ := newTestDeps(t)
	created, err := deps.Store.CreateProject("prod", []byte("tok"), []byte("nonce12byts"))
	if err != nil {
		t.Fatalf("CreateProject: %v", err)
	}
	// Pre-register one server with hcloud ID 1.
	if _, err := deps.Store.CreateServer(created.ID, store.Server{
		HCloudServerID: 1,
		Name:           "web-1",
		BaseServerType: "cpx11",
		TopServerType:  "cpx31",
		FallbackChain:  []string{"cpx31", "cpx11"},
		Mode:           "manual",
		Timezone:       "UTC",
	}); err != nil {
		t.Fatalf("CreateServer: %v", err)
	}

	stub := &fakeHetzner{
		servers: []*hetzner.Server{
			{ID: 1, Name: "web-1", ServerType: &hetzner.ServerType{Name: "cpx11"}},
			{ID: 2, Name: "web-2", ServerType: &hetzner.ServerType{Name: "cpx11"}},
		},
	}
	deps.APIFor = func(projectID int64) (hetzner.API, error) { return stub, nil }
	h := NewRouter(deps)

	req := authedRequest(t, http.MethodPost, "/api/projects/"+itoa(created.ID)+"/refresh", nil)
	rr := recorder(t, h, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("want 200, got %d (body=%q)", rr.Code, rr.Body.String())
	}

	var resp RefreshProjectResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(resp.Added) != 1 || resp.Added[0].HCloudServerID != 2 {
		t.Fatalf("want only server 2 added, got %+v", resp.Added)
	}
	if len(resp.Skipped) != 1 || resp.Skipped[0].HCloudServerID != 1 {
		t.Fatalf("want server 1 skipped, got %+v", resp.Skipped)
	}
}

// Ensure context import is referenced even if a future test doesn't use it.
var _ = context.Background
