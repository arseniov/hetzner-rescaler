package store

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// ErrNotFound is returned by Get* methods when no row matches.
var ErrNotFound = errors.New("store: not found")

// Temporary placeholder, replaced in Task 5.
type Server struct {
	ID             int64
	ProjectID      int64
	HCloudServerID int
	Name           string
	Mode           string
	BaseServerType string
	TopServerType  string
	FallbackChain  []string
	Timezone       string
}

// Project is a Hetzner Cloud project (identified by its API token).
type Project struct {
	ID                   int64
	Name                 string
	HCloudTokenEncrypted []byte
	HCloudTokenNonce     []byte
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

// CreateProject inserts a new project. The caller is responsible for
// encrypting the Hetzner token before calling.
func (s *Store) CreateProject(name string, tokenEnc, tokenNonce []byte) (*Project, error) {
	now := time.Now().UTC()
	res, err := s.db.Exec(
		`INSERT INTO projects (name, hcloud_token_encrypted, hcloud_token_nonce, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?)`,
		name, tokenEnc, tokenNonce, now.Unix(), now.Unix(),
	)
	if err != nil {
		return nil, fmt.Errorf("store: insert project: %w", err)
	}
	id, _ := res.LastInsertId()
	return &Project{
		ID:                   id,
		Name:                 name,
		HCloudTokenEncrypted: tokenEnc,
		HCloudTokenNonce:     tokenNonce,
		CreatedAt:            now,
		UpdatedAt:            now,
	}, nil
}

// GetProject returns a project by id.
func (s *Store) GetProject(id int64) (*Project, error) {
	row := s.db.QueryRow(
		`SELECT id, name, hcloud_token_encrypted, hcloud_token_nonce, created_at, updated_at
		 FROM projects WHERE id = ?`, id,
	)
	return scanProject(row)
}

// ListProjects returns all projects ordered by id.
func (s *Store) ListProjects() ([]*Project, error) {
	rows, err := s.db.Query(
		`SELECT id, name, hcloud_token_encrypted, hcloud_token_nonce, created_at, updated_at
		 FROM projects ORDER BY id`,
	)
	if err != nil {
		return nil, fmt.Errorf("store: list projects: %w", err)
	}
	defer rows.Close()

	var out []*Project
	for rows.Next() {
		p, err := scanProject(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

// DeleteProject removes a project. ON DELETE CASCADE removes its servers and
// their windows, events, and actions.
func (s *Store) DeleteProject(id int64) error {
	_, err := s.db.Exec(`DELETE FROM projects WHERE id = ?`, id)
	return err
}

// Temporary stubs for Task 4, replaced in Task 5.
// CreateServer inserts a new server under a project.
func (s *Store) CreateServer(projectID int64, srv Server) (*Server, error) {
	res, err := s.db.Exec(
		`INSERT INTO servers (project_id, hcloud_server_id, name, label, base_server_type, top_server_type, fallback_chain, mode, promote_state, timezone, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, NULL, ?, ?, ?)`,
		projectID, srv.HCloudServerID, srv.Name, "", srv.BaseServerType, srv.TopServerType,
		"[]", srv.Mode, srv.Timezone, time.Now().Unix(), time.Now().Unix(),
	)
	if err != nil {
		return nil, fmt.Errorf("store: insert server: %w", err)
	}
	id, _ := res.LastInsertId()
	srv.ID = id
	srv.ProjectID = projectID
	return &srv, nil
}

// GetServer returns a server by id.
func (s *Store) GetServer(id int64) (*Server, error) {
	row := s.db.QueryRow(
		`SELECT id, project_id, hcloud_server_id, name, base_server_type, top_server_type, mode, timezone
		 FROM servers WHERE id = ?`, id,
	)
	var srv Server
	err := row.Scan(&srv.ID, &srv.ProjectID, &srv.HCloudServerID, &srv.Name, &srv.BaseServerType, &srv.TopServerType, &srv.Mode, &srv.Timezone)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("store: scan server: %w", err)
	}
	return &srv, nil
}

// Add UNIQUE constraint on name via a separate migration step (idempotent).
// We add it here as a no-op ALTER that we ignore if it already exists.
func init() {
	// Best-effort UNIQUE on name. The store package's migrate() handles the
	// CREATE TABLE for projects, so we add the constraint via an index that
	// lives in the same migration. (No-op here; see Task 3's migration001.)
}

// rowScanner abstracts *sql.Row and *sql.Rows.
type rowScanner interface {
	Scan(dest ...any) error
}

func scanProject(r rowScanner) (*Project, error) {
	var (
		p       Project
		created int64
		updated int64
	)
	err := r.Scan(
		&p.ID, &p.Name, &p.HCloudTokenEncrypted, &p.HCloudTokenNonce,
		&created, &updated,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("store: scan project: %w", err)
	}
	p.CreatedAt = time.Unix(created, 0).UTC()
	p.UpdatedAt = time.Unix(updated, 0).UTC()
	return &p, nil
}