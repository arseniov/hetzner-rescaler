package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// User mirrors the columns of the `user` table that Better Auth's
// drizzle adapter creates. The Go engine does NOT create this table
// itself — the SPA's drizzle migrations do — so the columns here
// match what web/src/lib/server/schema.ts declares. We treat the
// table as read-only from the engine's side.
type User struct {
	ID            string
	Name          string
	Email         string
	EmailVerified bool
	Image         sql.NullString
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// Session mirrors the columns of the `session` table. Lookups by
// token go through GetSessionByToken.
type Session struct {
	ID        string
	Token     string
	UserID    string
	ExpiresAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	IPAddress sql.NullString
	UserAgent sql.NullString
}

// GetSessionByToken returns the session row matching token together
// with the owning user in a single round-trip via JOIN. Returns
// ErrNotFound if the token is unknown. Token expiry is the caller's
// responsibility — the session middleware checks `now < ExpiresAt`
// after loading so that "leaked tokens older than expiresAt but
// still in the table" can't sneak through.
func (s *Store) GetSessionByToken(token string) (*Session, *User, error) {
	return s.GetSessionByTokenContext(context.Background(), token)
}

// GetSessionByTokenContext is the context-aware variant used by the
// middleware so a slow DB query doesn't pin the request goroutine.
func (s *Store) GetSessionByTokenContext(ctx context.Context, token string) (*Session, *User, error) {
	if token == "" {
		return nil, nil, ErrNotFound
	}
	var (
		sess    Session
		user    User
		sessExp int64
		sessCr  int64
		sessUp  int64
		userCr  int64
		userUp  int64
	)
	row := s.db.QueryRowContext(ctx,
		`SELECT s.id, s.token, s.user_id, s.expires_at, s.created_at, s.updated_at, s.ip_address, s.user_agent,
		        u.id, u.name, u.email, u.email_verified, u.image, u.created_at, u.updated_at
		 FROM session s
		 JOIN user u ON u.id = s.user_id
		 WHERE s.token = ?`,
		token,
	)
	err := row.Scan(
		&sess.ID, &sess.Token, &sess.UserID, &sessExp, &sessCr, &sessUp, &sess.IPAddress, &sess.UserAgent,
		&user.ID, &user.Name, &user.Email, &user.EmailVerified, &user.Image, &userCr, &userUp,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil, ErrNotFound
	}
	if err != nil {
		return nil, nil, fmt.Errorf("store: scan session: %w", err)
	}
	sess.ExpiresAt = time.Unix(sessExp, 0).UTC()
	sess.CreatedAt = time.Unix(sessCr, 0).UTC()
	sess.UpdatedAt = time.Unix(sessUp, 0).UTC()
	user.CreatedAt = time.Unix(userCr, 0).UTC()
	user.UpdatedAt = time.Unix(userUp, 0).UTC()
	return &sess, &user, nil
}
