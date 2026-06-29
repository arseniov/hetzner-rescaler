package store

import (
	"database/sql"
	"fmt"
	"time"
)

// Event is one row in the append-only event log.
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

func (s *Store) ListAllEvents(limit int, serverID *int64, _ int) ([]*Event, error) {
	var (
		rows *sql.Rows
		err  error
	)
	if serverID != nil {
		if limit > 0 {
			rows, err = s.db.Query(
				`SELECT id, server_id, kind, from_type, to_type, started_at, finished_at, ok, error, triggered_by
				 FROM events WHERE server_id = ? ORDER BY id DESC LIMIT ?`,
				*serverID, limit,
			)
		} else {
			rows, err = s.db.Query(
				`SELECT id, server_id, kind, from_type, to_type, started_at, finished_at, ok, error, triggered_by
				 FROM events WHERE server_id = ? ORDER BY id DESC`,
				*serverID,
			)
		}
	} else {
		if limit > 0 {
			rows, err = s.db.Query(
				`SELECT id, server_id, kind, from_type, to_type, started_at, finished_at, ok, error, triggered_by
				 FROM events ORDER BY id DESC LIMIT ?`, limit,
			)
		} else {
			rows, err = s.db.Query(
				`SELECT id, server_id, kind, from_type, to_type, started_at, finished_at, ok, error, triggered_by
				 FROM events ORDER BY id DESC`,
			)
		}
	}
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
		)
		if err := rows.Scan(
			&e.ID, &e.ServerID, &e.Kind, &e.FromType, &e.ToType,
			&started, &finished, &okInt, &e.Error, &e.TriggeredBy,
		); err != nil {
			return nil, fmt.Errorf("store: scan event: %w", err)
		}
		e.StartedAt = time.Unix(started, 0).UTC()
		if finished.Valid {
			e.FinishedAt = time.Unix(finished.Int64, 0).UTC()
		}
		e.OK = okInt != 0
		out = append(out, &e)
	}
	return out, rows.Err()
}