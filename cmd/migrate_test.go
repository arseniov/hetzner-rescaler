package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jonamat/hetzner-rescaler/internal/store"
)

func TestMigrateImportsYAMLConfig(t *testing.T) {
	yaml := `hcloud_token: abc123
server_id: 15393230
base_server_name: cx11
top_server_name: cx21
hour_start: "09:00"
hour_stop: "20:00"
`
	yamlPath := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(yamlPath, []byte(yaml), 0600); err != nil {
		t.Fatalf("write yaml: %v", err)
	}

	dbPath := filepath.Join(t.TempDir(), "imported.db")
	st, err := store.Open(dbPath)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer st.Close()

	// Build a key in memory and pass directly
	key := make([]byte, 32)
	for i := range key {
		key[i] = byte(i)
	}
	t.Setenv("RESCALER_TOKEN_ENCRYPTION_KEY", "")

	if err := runMigrateImport(yamlPath, st, key); err != nil {
		t.Fatalf("runMigrateImport: %v", err)
	}

	projects, err := st.ListProjects()
	if err != nil {
		t.Fatalf("ListProjects: %v", err)
	}
	if len(projects) != 1 {
		t.Fatalf("got %d projects, want 1", len(projects))
	}

	servers, err := st.ListServersByProject(projects[0].ID)
	if err != nil {
		t.Fatalf("ListServersByProject: %v", err)
	}
	if len(servers) != 1 || servers[0].HCloudServerID != 15393230 {
		t.Fatalf("got %+v", servers)
	}
	if servers[0].Mode != "scheduled" {
		t.Fatalf("mode = %q, want scheduled", servers[0].Mode)
	}

	wins, err := st.ListWindows(servers[0].ID)
	if err != nil {
		t.Fatalf("ListWindows: %v", err)
	}
	if len(wins) != 1 {
		t.Fatalf("got %d windows, want 1 (created from hour_start/hour_stop)", len(wins))
	}
	if !strings.HasPrefix(wins[0].Label, "imported-") {
		t.Fatalf("window label = %q, want imported-...", wins[0].Label)
	}
}