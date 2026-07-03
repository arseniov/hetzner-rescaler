package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/jonamat/hetzner-rescaler/internal/crypto"
	"github.com/jonamat/hetzner-rescaler/internal/store"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

/* Migrate command */
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Import legacy YAML config into the SQLite database",
	Run:   runMigrate,
}

func init() {
	rootCmd.AddCommand(migrateCmd)
	migrateCmd.Flags().String("from", "", "Path to legacy YAML config (default: $RESCALER_YAML_PATH or ~/.hetzner-rescaler.yaml)")
}

func runMigrate(cmd *cobra.Command, args []string) {
	yamlPath, _ := cmd.Flags().GetString("from")
	if yamlPath == "" {
		yamlPath = os.Getenv("RESCALER_YAML_PATH")
	}
	if yamlPath == "" {
		fmt.Println("Usage: hetzner-rescaler migrate --from <path-to-yaml>")
		return
	}
	st, err := openStore()
	if err != nil {
		fmt.Println("open store:", err)
		return
	}
	defer st.Close()
	key, err := loadOrGenerateKey()
	if err != nil {
		fmt.Println("key:", err)
		return
	}
	if err := runMigrateImport(yamlPath, st, key); err != nil {
		fmt.Println("import:", err)
		return
	}
	fmt.Println("Import complete. Run `hetzner-rescaler status` to verify.")
}

type legacyConfig struct {
	HCloudToken    string `yaml:"hcloud_token"`
	ServerID       int    `yaml:"server_id"`
	BaseServerName string `yaml:"base_server_name"`
	TopServerName  string `yaml:"top_server_name"`
	HourStart      string `yaml:"hour_start"`
	HourStop       string `yaml:"hour_stop"`
}

// runMigrateImport reads the YAML at yamlPath and writes its contents into st.
func runMigrateImport(yamlPath string, st *store.Store, key []byte) error {
	data, err := os.ReadFile(yamlPath)
	if err != nil {
		return err
	}
	var cfg legacyConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return fmt.Errorf("parse yaml: %w", err)
	}
	if cfg.HCloudToken == "" || cfg.ServerID == 0 {
		return fmt.Errorf("yaml is missing required fields (hcloud_token, server_id)")
	}

	enc, nonce, err := crypto.Encrypt(key, []byte(cfg.HCloudToken))
	if err != nil {
		return err
	}
	proj, err := st.CreateProject("imported", enc, nonce)
	if err != nil {
		return err
	}

	// Default to UTC if TZ env var is unset so EvaluateWindows can load the location.
	timezone := os.Getenv("TZ")
	if timezone == "" {
		timezone = "UTC"
	}

	srv, err := st.CreateServer(proj.ID, store.Server{
		HCloudServerID: cfg.ServerID,
		Name:           fmt.Sprintf("server-%d", cfg.ServerID),
		BaseServerType: cfg.BaseServerName,
		TopServerType:  cfg.TopServerName,
		FallbackChain:  []string{cfg.TopServerName, cfg.BaseServerName},
		Mode:           "scheduled",
		Timezone:       timezone,
	})
	if err != nil {
		return err
	}
	if cfg.HourStart != "" && cfg.HourStop != "" {
		_, err := st.CreateWindow(srv.ID, store.Window{
			Label:      fmt.Sprintf("imported-%d", time.Now().Unix()),
			DaysOfWeek: 0b01111111, // all 7 days (Sun=0 ... Sat=6 under Go convention)
			StartTime:  cfg.HourStart,
			StopTime:   cfg.HourStop,
			TargetType: cfg.TopServerName,
			Enabled:    true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}