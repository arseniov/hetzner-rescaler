package api

import (
	"crypto/subtle"
	"encoding/json"
	"net/http"
)

// RequireInternalToken returns middleware that requires the
// X-Internal-Token request header to match expected using a constant-time
// compare. On mismatch or absence it writes 401 Unauthorized with a JSON
// body of {"error":"unauthorized"}.
//
// The expected token is supplied by NewRouter's caller (it comes from
// the RESCALER_INTERNAL_TOKEN environment variable in serve mode).
func RequireInternalToken(expected string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			got := r.Header.Get("X-Internal-Token")
			if subtle.ConstantTimeCompare([]byte(got), []byte(expected)) != 1 {
				writeJSONError(w, http.StatusUnauthorized, "unauthorized")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// writeJSONError writes a JSON error response. Used by the auth
// middleware and by handlers that need to surface error states.
func writeJSONError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": msg})
}