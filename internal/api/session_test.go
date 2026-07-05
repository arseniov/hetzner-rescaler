package api

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jonamat/hetzner-rescaler/internal/store"
)

// mustExec is a tiny helper so session tests can keep their boilerplate
// minimal when seeding fixture rows.
func mustExec(t *testing.T, st *store.Store, q string, args ...any) {
	t.Helper()
	if _, err := st.DB().Exec(q, args...); err != nil {
		t.Fatalf("exec %q: %v", q, err)
	}
}

// sign is the Go side of Better Auth's base64(HMAC-SHA256(secret, token))
// signing step. Tests use it to forge expected-good cookies without
// needing a live Better Auth instance.
func sign(secret, token string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(token))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

// withSession creates the session+user tables and inserts a row for
// the test to authenticate as.
func withSession(t *testing.T) (*store.Store, string, string) {
	t.Helper()
	st, err := store.OpenTemp()
	if err != nil {
		t.Fatalf("open temp store: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })

	// Better Auth's drizzle adapter creates these tables. We mirror
	// the schema here so GetSessionByTokenContext can SELECT against
	// them. The Go engine never owns these tables; the SPA's
	// migrations do.
	mustExec(t, st, `CREATE TABLE user (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		email_verified INTEGER NOT NULL DEFAULT 0,
		image TEXT,
		created_at INTEGER NOT NULL,
		updated_at INTEGER NOT NULL
	)`)
	mustExec(t, st, `CREATE TABLE session (
		id TEXT PRIMARY KEY,
		expires_at INTEGER NOT NULL,
		token TEXT NOT NULL UNIQUE,
		created_at INTEGER NOT NULL,
		updated_at INTEGER NOT NULL,
		ip_address TEXT,
		user_agent TEXT,
		user_id TEXT NOT NULL REFERENCES user(id) ON DELETE CASCADE
	)`)

	const token = "tok-abcdef0123456789"
	const secret = "test-secret"
	now := time.Now().UTC().Unix()
	mustExec(t, st, `INSERT INTO user (id, name, email, email_verified, created_at, updated_at) VALUES (?,?,?,?,?,?)`,
		"u1", "Alice", "alice@example.com", true, now, now)
	mustExec(t, st, `INSERT INTO session (id, expires_at, token, created_at, updated_at, ip_address, user_agent, user_id) VALUES (?,?,?,?,?,?,?,?)`,
		"s1", now+3600, token, now, now, "127.0.0.1", "vitest", "u1")
	return st, token, secret
}

func TestRequireAuth_AcceptsValidSessionCookie(t *testing.T) {
	st, token, secret := withSession(t)
	called, _ := newAuthOK(t)

	wrapped := RequireAuth("a-different-token", secret, st)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seen := AuthFromContext(r.Context())
		if seen == nil || seen.Session == nil || seen.User == nil {
			t.Fatalf("expected Session+User attached, got %#v", seen)
		}
		if seen.Session.Token != token {
			t.Fatalf("session.Token = %q, want %q", seen.Session.Token, token)
		}
		if seen.User.Email != "alice@example.com" {
			t.Fatalf("user.Email = %q", seen.User.Email)
		}
		*called = true
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/projects", nil)
	req.AddCookie(&http.Cookie{Name: "better-auth.session_token", Value: token + "." + sign(secret, token)})
	rr := httptest.NewRecorder()
	wrapped.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("want 200, got %d (body=%s)", rr.Code, rr.Body.String())
	}
	if !*called {
		t.Fatalf("downstream handler did not run")
	}
}

func TestRequireAuth_RejectsBadSignature(t *testing.T) {
	st, token, _ := withSession(t)
	called, next := newAuthOK(t)
	h := RequireAuth("", "the-real-secret", st)(next)

	req := httptest.NewRequest(http.MethodGet, "/api/projects", nil)
	// signature is computed with the WRONG secret.
	req.AddCookie(&http.Cookie{Name: "better-auth.session_token", Value: token + "." + sign("the-wrong-secret", token)})
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("want 401, got %d", rr.Code)
	}
	if *called {
		t.Fatalf("downstream handler should not run for bad signature")
	}
}

func TestRequireAuth_RejectsExpiredSession(t *testing.T) {
	st, _, secret := withSession(t)
	// Flip the row to expired.
	mustExec(t, st, `UPDATE session SET expires_at = ? WHERE token = ?`,
		time.Now().Add(-1*time.Hour).Unix(), "tok-abcdef0123456789")
	h := RequireAuth("", secret, st)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("downstream handler should not run for expired session")
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/projects", nil)
	req.AddCookie(&http.Cookie{Name: "better-auth.session_token",
		Value: "tok-abcdef0123456789." + sign(secret, "tok-abcdef0123456789")})
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("want 401, got %d", rr.Code)
	}
}

func TestRequireAuth_RejectsUnknownToken(t *testing.T) {
	st, _, secret := withSession(t)
	h := RequireAuth("", secret, st)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("downstream handler should not run for unknown token")
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/projects", nil)
	req.AddCookie(&http.Cookie{Name: "better-auth.session_token",
		Value: "ghost." + sign(secret, "ghost")})
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("want 401, got %d", rr.Code)
	}
}

func TestRequireAuth_PrefersInternalTokenOverCookie(t *testing.T) {
	st, token, secret := withSession(t)
	var seen *Auth
	wrapped := RequireAuth("machine-token", secret, st)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seen = AuthFromContext(r.Context())
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/projects", nil)
	// Both credentials present; the request must be admitted (token
	// matches) and the identity should reflect the InternalToken path
	// because we don't want to rely on a stored session row for M2M.
	req.Header.Set("X-Internal-Token", "machine-token")
	req.AddCookie(&http.Cookie{Name: "better-auth.session_token", Value: token + "." + sign(secret, token)})
	rr := httptest.NewRecorder()
	wrapped.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", rr.Code)
	}
	if seen == nil || !seen.InternalToken {
		t.Fatalf("expected InternalToken=true, got %#v", seen)
	}
}
