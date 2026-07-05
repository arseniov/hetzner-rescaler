package api

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// newAuthOK is a sentinel handler used to confirm a middleware
// admitted the request downstream.
func newAuthOK(t *testing.T) (called *bool, next http.Handler) {
	t.Helper()
	calledB := false
	return &calledB, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calledB = true
		w.WriteHeader(http.StatusOK)
	})
}

func TestRequireAuth_RejectsMissing(t *testing.T) {
	called, next := newAuthOK(t)
	h := RequireAuth("secret-abc", "", nil)(next)
	req := httptest.NewRequest(http.MethodGet, "/api/projects", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("want 401, got %d", rr.Code)
	}
	if *called {
		t.Fatalf("downstream handler should not run")
	}
	var body map[string]string
	_ = json.Unmarshal(rr.Body.Bytes(), &body)
	if body["error"] != "unauthorized" {
		t.Fatalf("want error=unauthorized, got %q", body["error"])
	}
}

func TestRequireAuth_RejectsWrongToken(t *testing.T) {
	called, next := newAuthOK(t)
	h := RequireAuth("secret-abc", "", nil)(next)
	req := httptest.NewRequest(http.MethodGet, "/api/projects", nil)
	req.Header.Set("X-Internal-Token", "wrong")
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("want 401, got %d", rr.Code)
	}
	if *called {
		t.Fatalf("downstream handler should not run")
	}
}

func TestRequireAuth_AllowsCorrectToken(t *testing.T) {
	called, next := newAuthOK(t)
	h := RequireAuth("secret-abc", "", nil)(next)
	req := httptest.NewRequest(http.MethodGet, "/api/projects", nil)
	req.Header.Set("X-Internal-Token", "secret-abc")
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", rr.Code)
	}
	if !*called {
		t.Fatalf("downstream handler not invoked")
	}
}

func TestRequireAuth_UsesConstantTimeCompare(t *testing.T) {
	// Two tokens of equal length but different content should be
	// rejected and the body should be the JSON error envelope.
	called, next := newAuthOK(t)
	h := RequireAuth("aaaaaaaaaaaaaaaa", "", nil)(next)
	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	req.Header.Set("X-Internal-Token", "bbbbbbbbbbbbbbbb")
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("want 401, got %d", rr.Code)
	}
	if got, _ := io.ReadAll(rr.Body); len(got) == 0 {
		t.Fatalf("expected JSON error body, got empty")
	}
	if *called {
		t.Fatalf("downstream handler should not run for mismatched token")
	}
}

func TestRequireAuth_NoSessionSecretFallsBackToTokenOnly(t *testing.T) {
	// When SessionSecret is empty, the cookie path is skipped
	// entirely — even a perfectly valid cookie should not bypass
	// the missing X-Internal-Token header.
	h := RequireAuth("secret-abc", "", nil)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("downstream handler should not run")
	}))
	req := httptest.NewRequest(http.MethodGet, "/api/projects", nil)
	req.AddCookie(&http.Cookie{Name: "better-auth.session_token", Value: "abc.def"})
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("want 401 with empty SessionSecret, got %d", rr.Code)
	}
}

func TestAuthFromContext_AttachesIdentity(t *testing.T) {
	var seen *Auth
	wrapped := RequireAuth("secret-abc", "", nil)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seen = AuthFromContext(r.Context())
		w.WriteHeader(http.StatusOK)
	}))
	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	req.Header.Set("X-Internal-Token", "secret-abc")
	rr := httptest.NewRecorder()
	wrapped.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", rr.Code)
	}
	if seen == nil {
		t.Fatalf("AuthFromContext returned nil")
	}
	if !seen.InternalToken {
		t.Fatalf("expected InternalToken=true")
	}
	if seen.Session != nil {
		t.Fatalf("expected Session=nil for token-only auth")
	}
}
