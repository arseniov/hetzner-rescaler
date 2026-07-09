package store

import (
	"database/sql"
	"path/filepath"
	"testing"
)

func TestMigrationAddsPhaseColumn(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "test.db")
	s, _ := Open(dbPath)
	defer s.Close()

	rows, err := s.DB().Query(`PRAGMA table_info(events)`)
	if err != nil {
		t.Fatalf("PRAGMA: %v", err)
	}
	defer rows.Close()
	var hasPhase bool
	for rows.Next() {
		var (
			cid            int
			name, ctype    string
			notnull, pk    int
			dfltValue      sql.NullString
		)
		if err := rows.Scan(&cid, &name, &ctype, &notnull, &dfltValue, &pk); err != nil {
			t.Fatalf("scan: %v", err)
		}
		if name == "phase" {
			hasPhase = true
			if ctype != "TEXT" {
				t.Fatalf("phase column type = %q, want TEXT", ctype)
			}
		}
	}
	if !hasPhase {
		t.Fatalf("events table missing 'phase' column")
	}
}

func TestSchemaVersionIs3(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "test.db")
	s, _ := Open(dbPath)
	defer s.Close()
	if v := readSchemaVersion(t, s); v != 3 {
		t.Fatalf("schema version = %d, want 3", v)
	}
}

func TestMigrationsAreForwardOnly(t *testing.T) {
	// Apply migrations twice in a row, version must not regress.
	dbPath := filepath.Join(t.TempDir(), "test.db")
	s, _ := Open(dbPath)
	defer s.Close()
	v1 := readSchemaVersion(t, s)
	s.Close()

	s2, _ := Open(dbPath)
	defer s2.Close()
	v2 := readSchemaVersion(t, s2)

	if v2 < v1 {
		t.Fatalf("schema version regressed: %d -> %d", v1, v2)
	}
}

func TestRequiredTablesExist(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "test.db")
	s, _ := Open(dbPath)
	defer s.Close()

	want := []string{"projects", "servers", "windows", "events", "actions", "schema_version"}
	for _, table := range want {
		var n int
		err := s.DB().QueryRow(
			`SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name=?`, table,
		).Scan(&n)
		if err != nil {
			t.Fatalf("query sqlite_master for %s: %v", table, err)
		}
		if n != 1 {
			t.Fatalf("table %q not found (count=%d)", table, n)
		}
	}
}

func readSchemaVersion(t *testing.T, s *Store) int {
	t.Helper()
	var v int
	if err := s.DB().QueryRow(`SELECT MAX(version) FROM schema_version`).Scan(&v); err != nil {
		t.Fatalf("read schema_version: %v", err)
	}
	return v
}
