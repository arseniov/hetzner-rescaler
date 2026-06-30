package web

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSPAHandler_ServesIndexHTML(t *testing.T) {
	rec := httptest.NewRecorder()
	SPAHandler().ServeHTTP(rec, httptest.NewRequest("GET", "/", nil))

	if rec.Code != 200 {
		t.Fatalf("want 200, got %d", rec.Code)
	}
	body, _ := io.ReadAll(rec.Body)
	if !strings.Contains(string(body), "<html") {
		t.Fatalf("expected HTML response, got %q", string(body))
	}
}

func TestSPAHandler_FallbackForArbitraryPath(t *testing.T) {
	// /projects/42 is not a real file in the build output; it should
	// fall back to index.html so the SPA's client-side router can take
	// over.
	rec := httptest.NewRecorder()
	SPAHandler().ServeHTTP(rec, httptest.NewRequest("GET", "/projects/42", nil))
	if rec.Code != 200 {
		t.Fatalf("want 200, got %d", rec.Code)
	}
	body, _ := io.ReadAll(rec.Body)
	if !strings.Contains(string(body), "<html") {
		t.Fatalf("expected fallback to index.html, got %q", string(body))
	}
}