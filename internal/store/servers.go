package store

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// Server is a registered Hetzner Cloud server with its rescale configuration.
type Server struct {
	ID             int64
	ProjectID      int64
	HCloudServerID int
	Name           string
	Label          string
	BaseServerType string
	TopServerType  string
	FallbackChain  []string // stored as JSON
	Mode           string   // 'manual' | 'auto_promote' | 'scheduled'
	PromoteState   *string  // nil unless mode == 'auto_promote'
	Timezone       string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	store          *Store
}

// Store returns the *Store that owns this server. It is set on rows created
// via CreateServer but is nil on rows loaded directly from SQL by scanServer.
func (srv *Server) Store() *Store { return srv.store }

// Window is a scheduled scale-up window for a server in 'scheduled' mode.
type Window struct {
	ID         int64
	ServerID   int64
	Label      string
	DaysOfWeek int    // bitmask: bit 0=Mon, bit 6=Sun
	StartTime  string // "HH:MM"
	StopTime   string // "HH:MM"
	TargetType string
	Enabled    bool
}

func (s *Store) CreateServer(projectID int64, srv Server) (*Server, error) {
	chain, err := json.Marshal(srv.FallbackChain)
	if err != nil {
		return nil, fmt.Errorf("store: marshal fallback_chain: %w", err)
	}
	now := time.Now().UTC()
	res, err := s.db.Exec(
		`INSERT INTO servers
		 (project_id, hcloud_server_id, name, label, base_server_type, top_server_type,
		  fallback_chain, mode, promote_state, timezone, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		projectID, srv.HCloudServerID, srv.Name, srv.Label,
		srv.BaseServerType, srv.TopServerType, string(chain), srv.Mode,
		srv.PromoteState, srv.Timezone, now.Unix(), now.Unix(),
	)
	if err != nil {
		return nil, fmt.Errorf("store: insert server: %w", err)
	}
	id, _ := res.LastInsertId()
	srv.ID = id
	srv.ProjectID = projectID
	srv.CreatedAt = now
	srv.UpdatedAt = now
	srv.store = s
	return &srv, nil
}

func (s *Store) GetServer(id int64) (*Server, error) {
	row := s.db.QueryRow(`SELECT ` + serverCols + ` FROM servers WHERE id = ?`, id)
	return scanServer(row)
}

func (s *Store) ListServersByProject(projectID int64) ([]*Server, error) {
	rows, err := s.db.Query(`SELECT ` + serverCols + ` FROM servers WHERE project_id = ? ORDER BY id`, projectID)
	if err != nil {
		return nil, fmt.Errorf("store: list servers: %w", err)
	}
	defer rows.Close()
	var out []*Server
	for rows.Next() {
		srv, err := scanServer(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, srv)
	}
	return out, rows.Err()
}

func (s *Store) ListAllServers() ([]*Server, error) {
	rows, err := s.db.Query(`SELECT ` + serverCols + ` FROM servers ORDER BY id`)
	if err != nil {
		return nil, fmt.Errorf("store: list all servers: %w", err)
	}
	defer rows.Close()
	var out []*Server
	for rows.Next() {
		srv, err := scanServer(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, srv)
	}
	return out, rows.Err()
}

func (s *Store) UpdateServer(srv Server) error {
	chain, err := json.Marshal(srv.FallbackChain)
	if err != nil {
		return fmt.Errorf("store: marshal fallback_chain: %w", err)
	}
	now := time.Now().UTC()
	res, err := s.db.Exec(
		`UPDATE servers SET
		   name=?, label=?, base_server_type=?, top_server_type=?,
		   fallback_chain=?, mode=?, promote_state=?, timezone=?, updated_at=?
		 WHERE id=?`,
		srv.Name, srv.Label, srv.BaseServerType, srv.TopServerType,
		string(chain), srv.Mode, srv.PromoteState, srv.Timezone, now.Unix(), srv.ID,
	)
	if err != nil {
		return fmt.Errorf("store: update server: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *Store) DeleteServer(id int64) error {
	_, err := s.db.Exec(`DELETE FROM servers WHERE id = ?`, id)
	return err
}

const serverCols = `id, project_id, hcloud_server_id, name, label, base_server_type, top_server_type,
	fallback_chain, mode, promote_state, timezone, created_at, updated_at`

func scanServer(r rowScanner) (*Server, error) {
	var (
		srv      Server
		chainRaw string
		promote  sql.NullString
		created  int64
		updated  int64
	)
	err := r.Scan(
		&srv.ID, &srv.ProjectID, &srv.HCloudServerID, &srv.Name, &srv.Label,
		&srv.BaseServerType, &srv.TopServerType, &chainRaw, &srv.Mode, &promote,
		&srv.Timezone, &created, &updated,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("store: scan server: %w", err)
	}
	if err := json.Unmarshal([]byte(chainRaw), &srv.FallbackChain); err != nil {
		return nil, fmt.Errorf("store: unmarshal fallback_chain: %w", err)
	}
	if promote.Valid {
		v := promote.String
		srv.PromoteState = &v
	}
	srv.CreatedAt = time.Unix(created, 0).UTC()
	srv.UpdatedAt = time.Unix(updated, 0).UTC()
	return &srv, nil
}

func (s *Store) CreateWindow(serverID int64, w Window) (*Window, error) {
	enabled := 0
	if w.Enabled {
		enabled = 1
	}
	res, err := s.db.Exec(
		`INSERT INTO windows (server_id, label, days_of_week, start_time, stop_time, target_type, enabled)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		serverID, w.Label, w.DaysOfWeek, w.StartTime, w.StopTime, w.TargetType, enabled,
	)
	if err != nil {
		return nil, fmt.Errorf("store: insert window: %w", err)
	}
	id, _ := res.LastInsertId()
	w.ID = id
	w.ServerID = serverID
	return &w, nil
}

func (s *Store) ListWindows(serverID int64) ([]*Window, error) {
	rows, err := s.db.Query(
		`SELECT id, server_id, label, days_of_week, start_time, stop_time, target_type, enabled
		 FROM windows WHERE server_id = ? ORDER BY id`, serverID,
	)
	if err != nil {
		return nil, fmt.Errorf("store: list windows: %w", err)
	}
	defer rows.Close()
	var out []*Window
	for rows.Next() {
		var w Window
		var enabled int
		if err := rows.Scan(&w.ID, &w.ServerID, &w.Label, &w.DaysOfWeek, &w.StartTime, &w.StopTime, &w.TargetType, &enabled); err != nil {
			return nil, fmt.Errorf("store: scan window: %w", err)
		}
		w.Enabled = enabled != 0
		out = append(out, &w)
	}
	return out, rows.Err()
}

func (s *Store) DeleteWindow(id int64) error {
	_, err := s.db.Exec(`DELETE FROM windows WHERE id = ?`, id)
	return err
}