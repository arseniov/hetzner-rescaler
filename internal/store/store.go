// Package store is the only package in the engine that knows SQL.
// All other internal/* packages consume typed structs.
package store

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/jonamat/hetzner-rescaler/internal/broadcast"

	_ "modernc.org/sqlite" // pure-Go sqlite driver
)

// currentSchemaVersion is incremented whenever migrations.go changes.
const currentSchemaVersion = 3

// Store is the SQLite-backed persistence layer.
type Store struct {
	db  *sql.DB
	hub *broadcast.Hub[Event]
}

// SetBroadcastHub attaches an in-process pub/sub hub that receives every
// event after it is successfully inserted. Pass nil to detach.
func (s *Store) SetBroadcastHub(hub *broadcast.Hub[Event]) {
	s.hub = hub
}

// EventHub returns the broadcast hub attached via SetBroadcastHub, or nil
// if none has been attached. SSE handlers and other in-process subscribers
// use this to receive live events as they are inserted.
func (s *Store) EventHub() *broadcast.Hub[Event] {
	return s.hub
}

// Open opens (or creates) the database at path and runs all migrations.
func Open(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path+"?_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)")
	if err != nil {
		return nil, fmt.Errorf("store: open: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("store: ping: %w", err)
	}
	s := &Store{db: db}
	if err := s.migrate(); err != nil {
		s.Close()
		return nil, fmt.Errorf("store: migrate: %w", err)
	}
	return s, nil
}

// DB returns the underlying *sql.DB. Use sparingly — prefer typed CRUD methods.
func (s *Store) DB() *sql.DB { return s.db }

// Close closes the database.
func (s *Store) Close() error { return s.db.Close() }

// OpenTemp opens a database in a temp file. The caller must Close it.
func OpenTemp() (*Store, error) {
	f, err := os.CreateTemp("", "rescaler-test-*.db")
	if err != nil {
		return nil, err
	}
	f.Close()
	return Open(f.Name())
}
