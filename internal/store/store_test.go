package store

import (
	"path/filepath"
	"testing"
)

func TestOpenCreatesDBAndRunsMigrations(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "test.db")
	s, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer s.Close()

	if s.DB() == nil {
		t.Fatal("DB() returned nil")
	}

	// schema_version should be at the current version
	var v int
	if err := s.DB().QueryRow(`SELECT version FROM schema_version LIMIT 1`).Scan(&v); err != nil {
		t.Fatalf("query schema_version: %v", err)
	}
	if v != currentSchemaVersion {
		t.Fatalf("schema_version = %d, want %d", v, currentSchemaVersion)
	}
}

func TestOpenIsIdempotent(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "test.db")
	s1, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open #1: %v", err)
	}
	s1.Close()

	s2, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open #2: %v", err)
	}
	defer s2.Close()
}
