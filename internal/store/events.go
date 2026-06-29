package store

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

// Event is one row in the append-only event log.
//
// FinishedAt is stored as a value (not a pointer) for ergonomics: callers
// leave it at the Go zero value when the rescale is still in progress and
// treat FinishedAt.IsZero() as the "pending" sentinel. AppendEvent handles
// the conversion to a nullable INTEGER on the way to SQLite, and scanEvents
// leaves FinishedAt at the zero value when the column is NULL.
type Event struct {
	ID          int64
	ServerID    int64
	Kind        string
	FromType    string
	ToType      string
	StartedAt   time.Time
	FinishedAt  time.Time
	OK          bool
	Error       string
	TriggeredBy string
}

func (s *Store) AppendEvent(e Event) (int64, error) {
	var finished *int64
	if !e.FinishedAt.IsZero() {
		ts := e.FinishedAt.Unix()
		finished = &ts
	}
	okInt := 0
	if e.OK {
		okInt = 1
	}
	res, err := s.db.Exec(
		`INSERT INTO events (server_id, kind, from_type, to_type, started_at, finished_at, ok, error, triggered_by)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		e.ServerID, e.Kind, e.FromType, e.ToType,
		e.StartedAt.Unix(), finished, okInt, e.Error, e.TriggeredBy,
	)
	if err != nil {
		return 0, fmt.Errorf("store: insert event: %w", err)
	}
	id, _ := res.LastInsertId()
	return id, nil
}

func (s *Store) ListEventsByServer(serverID int64, limit int) ([]*Event, error) {
	rows, err := s.db.Query(
		`SELECT id, server_id, kind, from_type, to_type, started_at, finished_at, ok, error, triggered_by
		 FROM events WHERE server_id = ? ORDER BY id DESC LIMIT ?`,
		serverID, limit,
	)
	if err != nil {
		return nil, fmt.Errorf("store: list events: %w", err)
	}
	defer rows.Close()
	return scanEvents(rows)
}

// ListAllEvents returns events across all servers, newest first. If serverID
// is non-nil, only events for that server are returned. If limit > 0, at most
// limit rows are returned; otherwise all matching rows are returned.
func (s *Store) ListAllEvents(limit int, serverID *int64) ([]*Event, error) {
	var (
		clauses []string
		args    []any
	)
	if serverID != nil {
		clauses = append(clauses, "server_id = ?")
		args = append(args, *serverID)
	}
	q := `SELECT id, server_id, kind, from_type, to_type, started_at, finished_at, ok, error, triggered_by FROM events`
	if len(clauses) > 0 {
		q += " WHERE " + strings.Join(clauses, " AND ")
	}
	q += " ORDER BY id DESC"
	if limit > 0 {
		q += " LIMIT ?"
		args = append(args, limit)
	}
	rows, err := s.db.Query(q, args...)
	if err != nil {
		return nil, fmt.Errorf("store: list all events: %w", err)
	}
	defer rows.Close()
	return scanEvents(rows)
}

func scanEvents(rows *sql.Rows) ([]*Event, error) {
	var out []*Event
	for rows.Next() {
		var (
			e        Event
			started  int64
			finished sql.NullInt64
			okInt    int
			fromType sql.NullString
			toType   sql.NullString
			errMsg   sql.NullString
		)
		if err := rows.Scan(
			&e.ID, &e.ServerID, &e.Kind, &fromType, &toType,
			&started, &finished, &okInt, &errMsg, &e.TriggeredBy,
		); err != nil {
			return nil, fmt.Errorf("store: scan event: %w", err)
		}
		e.FromType = fromType.String
		e.ToType = toType.String
		e.Error = errMsg.String
		e.StartedAt = time.Unix(started, 0).UTC()
		if finished.Valid {
			e.FinishedAt = time.Unix(finished.Int64, 0).UTC()
		}
		e.OK = okInt != 0
		out = append(out, &e)
	}
	return out, rows.Err()
}
