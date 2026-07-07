package store

import "fmt"

// migrations is the ordered list of schema migrations. Append-only.
// Each migration is wrapped in a transaction. Migrations must be idempotent
// (safe to re-run on an already-current DB) so that crash recovery is safe.
var migrations = []func(*Store) error{
	migration001_initial,
	migration002_add_phase,
}

func (s *Store) migrate() error {
	if _, err := s.db.Exec(`CREATE TABLE IF NOT EXISTS schema_version (version INTEGER PRIMARY KEY)`); err != nil {
		return fmt.Errorf("create schema_version: %w", err)
	}
	for i, fn := range migrations {
		step := i + 1
		if step <= s.currentVersion() {
			continue
		}
		tx, err := s.db.Begin()
		if err != nil {
			return fmt.Errorf("begin migration %d: %w", step, err)
		}
		if err := fn(s); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("apply migration %d: %w", step, err)
		}
		if _, err := tx.Exec(`INSERT INTO schema_version(version) VALUES (?)`, step); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("record schema_version %d: %w", step, err)
		}
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("commit migration %d: %w", step, err)
		}
	}
	return nil
}

func (s *Store) currentVersion() int {
	var v int
	row := s.db.QueryRow(`SELECT COALESCE(MAX(version), 0) FROM schema_version`)
	if err := row.Scan(&v); err != nil {
		return 0
	}
	return v
}

func migration001_initial(s *Store) error {
	stmts := []string{
		`CREATE TABLE projects (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE,
			hcloud_token_encrypted BLOB NOT NULL,
			hcloud_token_nonce BLOB NOT NULL,
			created_at INTEGER NOT NULL,
			updated_at INTEGER NOT NULL
		)`,
		`CREATE TABLE servers (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			project_id INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
			hcloud_server_id INTEGER NOT NULL,
			name TEXT NOT NULL,
			label TEXT NOT NULL DEFAULT '',
			base_server_type TEXT NOT NULL,
			top_server_type TEXT NOT NULL,
			fallback_chain TEXT NOT NULL DEFAULT '[]',
			mode TEXT NOT NULL CHECK(mode IN ('manual','auto_promote','scheduled')),
			promote_state TEXT,
			timezone TEXT NOT NULL DEFAULT 'UTC',
			created_at INTEGER NOT NULL,
			updated_at INTEGER NOT NULL,
			UNIQUE(project_id, hcloud_server_id)
		)`,
		`CREATE TABLE windows (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			server_id INTEGER NOT NULL REFERENCES servers(id) ON DELETE CASCADE,
			label TEXT NOT NULL,
			days_of_week INTEGER NOT NULL,
			start_time TEXT NOT NULL,
			stop_time TEXT NOT NULL,
			target_type TEXT NOT NULL,
			enabled INTEGER NOT NULL DEFAULT 1,
			UNIQUE(server_id, label)
		)`,
		`CREATE TABLE events (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			server_id INTEGER NOT NULL REFERENCES servers(id) ON DELETE CASCADE,
			kind TEXT NOT NULL,
			from_type TEXT,
			to_type TEXT,
			started_at INTEGER NOT NULL,
			finished_at INTEGER,
			ok INTEGER NOT NULL,
			error TEXT,
			triggered_by TEXT NOT NULL
		)`,
		`CREATE INDEX idx_events_server_id ON events(server_id)`,
		`CREATE TABLE actions (
			server_id INTEGER PRIMARY KEY REFERENCES servers(id) ON DELETE CASCADE,
			kind TEXT NOT NULL,
			started_at INTEGER NOT NULL,
			locked_until INTEGER NOT NULL
		)`,
	}
	for _, q := range stmts {
		if _, err := s.db.Exec(q); err != nil {
			return fmt.Errorf("exec: %s: %w", q[:40], err)
		}
	}
	return nil
}

func migration002_add_phase(s *Store) error {
	if _, err := s.db.Exec(`ALTER TABLE events ADD COLUMN phase TEXT`); err != nil {
		return fmt.Errorf("add phase column: %w", err)
	}
	return nil
}
