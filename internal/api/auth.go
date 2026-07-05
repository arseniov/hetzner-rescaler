package api

import (
	"context"
	"crypto/subtle"
	"encoding/json"
	"net/http"

	"github.com/jonamat/hetzner-rescaler/internal/store"
)

// authContextKey is the request-scoped context key under which the
// middleware stores the authenticated identity (if any). Handlers
// that need to know whether the caller is a human (session) or a
// machine (InternalToken) can pluck it out via AuthFromContext.
type authContextKey struct{}

// Auth is what the middleware attaches to the request context so
// downstream handlers can introspect who made the call. At least
// one of InternalToken or Session is always non-nil when the
// middleware admitted the request.
type Auth struct {
	// InternalToken is true when the request was admitted via the
	// X-Internal-Token header. This is the path CLI scripts and
	// EventSource-with-?token= use.
	InternalToken bool

	// Session is the row from Better Auth's `session` table when the
	// request was admitted via a verified session cookie. nil for
	// InternalToken-only requests.
	Session *store.Session

	// User is the session's owning user. nil for InternalToken-only.
	User *store.User
}

// internalTokenMatched reports whether the X-Internal-Token header on
// req matches expected. The constant-time compare keeps the check
// immune to timing attacks.
func internalTokenMatched(req *http.Request, expected string) bool {
	if expected == "" {
		return false
	}
	return subtle.ConstantTimeCompare(
		[]byte(req.Header.Get("X-Internal-Token")),
		[]byte(expected),
	) == 1
}

// RequireAuth admits the request if EITHER the X-Internal-Token
// header matches (machine-to-machine / CLI / EventSource ?token=
// fallback) OR a valid Better Auth session cookie is present
// (browser SPA). This is defense in depth: the SPA already sends
// the internal token (baked into its bundle), but accepting either
// means a future rotation of RESCALER_INTERNAL_TOKEN does not need
// to wait for every CLI client to update before users regain access.
// On rejection, writes 401 with {"error":"unauthorized"}.
//
// sessionSecret may be empty — in that case RequireAuth degrades to
// InternalToken-only behaviour, matching the original middleware.
func RequireAuth(internalToken, sessionSecret string, st *store.Store) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			a := &Auth{}
			if internalTokenMatched(r, internalToken) {
				a.InternalToken = true
			} else if sessionSecret != "" && st != nil {
				sess, user, err := verifySessionCookie(r, sessionSecret, st)
				if err == nil {
					a.Session = sess
					a.User = user
				}
			}
			if !a.InternalToken && a.Session == nil {
				writeJSONError(w, http.StatusUnauthorized, "unauthorized")
				return
			}
			ctx := context.WithValue(r.Context(), authContextKey{}, a)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// AuthFromContext returns the Auth attached by RequireAuth, or nil
// if the request was not admitted via RequireAuth (e.g. for
// /api/healthz which has no auth wrapping).
func AuthFromContext(ctx context.Context) *Auth {
	v, _ := ctx.Value(authContextKey{}).(*Auth)
	return v
}

// InternalTokenMatched is exported for use by the SSE middleware,
// which lifts a ?token= query parameter into a header before calling
// the standard check.
func InternalTokenMatched(req *http.Request, expected string) bool {
	return internalTokenMatched(req, expected)
}

// writeJSONError writes a JSON error response. Used by the auth
// middleware and by handlers that need to surface error states.
func writeJSONError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
