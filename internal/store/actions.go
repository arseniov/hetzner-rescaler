package store

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// ErrLocked is returned by AcquireAction when an action is already in flight.
var ErrLocked = errors.New("store: action already in flight")

// AcquireAction attempts to lock the per-server action slot for `ttl`.
// Returns (true, nil) if acquired; (false, ErrLocked) if another action is in
// flight and not yet expired; (false, nil) and overwrites the lock if the
// existing lock is past its `locked_until`.
func (s *Store) AcquireAction(serverID int64, kind string, ttl time.Duration) (bool, error) {
	now := time.Now().UTC()
	until := now.Add(ttl)
	untilTs := until.Unix()

	// Check existing
	var (
		existingUntil int64
		existingKind  string
	)
	row := s.db.QueryRow(`SELECT kind, locked_until FROM actions WHERE server_id = ?`, serverID)
	err := row.Scan(&existingKind, &existingUntil)
	if errors.Is(err, sql.ErrNoRows) {
		// Free, insert
		_, err := s.db.Exec(
			`INSERT INTO actions (server_id, kind, started_at, locked_until) VALUES (?, ?, ?, ?)`,
			serverID, kind, now.Unix(), untilTs,
		)
		if err != nil {
			return false, fmt.Errorf("store: insert action: %w", err)
		}
		return true, nil
	}
	if err != nil {
		return false, fmt.Errorf("store: query action: %w", err)
	}
	if existingUntil > now.Unix() {
		return false, ErrLocked
	}
	// Stale: take over
	_, err = s.db.Exec(
		`UPDATE actions SET kind = ?, started_at = ?, locked_until = ? WHERE server_id = ?`,
		kind, now.Unix(), untilTs, serverID,
	)
	if err != nil {
		return false, fmt.Errorf("store: update action: %w", err)
	}
	return false, nil
}

// ReleaseAction releases the per-server action lock. No-op if not held.
func (s *Store) ReleaseAction(serverID int64) error {
	_, err := s.db.Exec(`DELETE FROM actions WHERE server_id = ?`, serverID)
	return err
}