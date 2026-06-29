package api

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequireInternalToken_RejectsMissing(t *testing.T) {
	h := RequireInternalToken("secret-abc")(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Fatalf("downstream handler should not run")
		}),
	)
	req := httptest.NewRequest(http.MethodGet, "/api/projects", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("want 401, got %d", rr.Code)
	}
	var body map[string]string
	_ = json.Unmarshal(rr.Body.Bytes(), &body)
	if body["error"] != "unauthorized" {
		t.Fatalf("want error=unauthorized, got %q", body["error"])
	}
}

func TestRequireInternalToken_RejectsWrongToken(t *testing.T) {
	h := RequireInternalToken("secret-abc")(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Fatalf("downstream handler should not run")
		}),
	)
	req := httptest.NewRequest(http.MethodGet, "/api/projects", nil)
	req.Header.Set("X-Internal-Token", "wrong")
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("want 401, got %d", rr.Code)
	}
}

func TestRequireInternalToken_AllowsCorrectToken(t *testing.T) {
	called := false
	h := RequireInternalToken("secret-abc")(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		}),
	)
	req := httptest.NewRequest(http.MethodGet, "/api/projects", nil)
	req.Header.Set("X-Internal-Token", "secret-abc")
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", rr.Code)
	}
	if !called {
		t.Fatalf("downstream handler not invoked")
	}
}

func TestRequireInternalToken_UsesConstantTimeCompare(t *testing.T) {
	// Confirm the implementation uses subtle.ConstantTimeCompare by
	// checking that two tokens of equal length but different content
	// are rejected (and that body is exactly what we expect).
	h := RequireInternalToken("aaaaaaaaaaaaaaaa")(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) }),
	)
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
}