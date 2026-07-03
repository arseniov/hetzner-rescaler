package store

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// ErrNotFound is returned by Get* methods when no row matches.
var ErrNotFound = errors.New("store: not found")

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