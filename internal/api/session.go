package api

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"strings"
	"time"

	"github.com/jonamat/hetzner-rescaler/internal/store"
)

// Better Auth's compact / legacy base64-HMAC cookie format:
//
//	<token>.<signature>
//
// where token is the session row's `token` column (32 ASCII chars in
// our case) and signature is base64(HMAC-SHA256(secret, token)).
// The signature is PADDED base64 — not URL-safe — as produced by the
// browser-native btoa() in Better Auth's makeSignature. We verify it
// here against the same secret so an attacker without the secret
// cannot forge a session row's token even if they could guess one.

const (
	betterAuthCookieName = "better-auth.session_token"
)

// verifySessionCookie parses the better-auth.session_token cookie
// from req, validates its HMAC against secret, looks up the session
// row in store, and confirms the session hasn't expired. Returns
// the session + user on success. Returns errInvalidSession for any
// failure mode (no cookie, malformed, bad signature, unknown token,
// expired session) so callers can treat all rejection paths the
// same.
func verifySessionCookie(req *http.Request, secret string, st *store.Store) (*store.Session, *store.User, error) {
	if secret == "" || st == nil {
		return nil, nil, errInvalidSession
	}
	c, err := req.Cookie(betterAuthCookieName)
	if err != nil || c.Value == "" {
		return nil, nil, errInvalidSession
	}
	parts := strings.SplitN(c.Value, ".", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return nil, nil, errInvalidSession
	}
	token, sigB64 := parts[0], parts[1]

	// HMAC-SHA256(secret, token) → base64 (padded, NOT url-safe).
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(token))
	expectedSig, err := base64.StdEncoding.DecodeString(sigB64)
	if err != nil {
		return nil, nil, errInvalidSession
	}
	if !hmac.Equal(expectedSig, mac.Sum(nil)) {
		return nil, nil, errInvalidSession
	}

	// Cookie signature is fine — now consult the DB. This catches
	// revoked sessions (admin sign-out, sessions table truncated,
	// cookie replay after expiry).
	ctx, cancel := context.WithTimeout(req.Context(), 1*time.Second)
	defer cancel()
	sess, user, err := st.GetSessionByTokenContext(ctx, token)
	if err != nil {
		return nil, nil, errInvalidSession
	}
	if !sess.ExpiresAt.IsZero() && sess.ExpiresAt.Before(time.Now().UTC()) {
		return nil, nil, errInvalidSession
	}
	return sess, user, nil
}

var errInvalidSession = sentinelInvalid{}

type sentinelInvalid struct{}

func (sentinelInvalid) Error() string { return "invalid session" }
