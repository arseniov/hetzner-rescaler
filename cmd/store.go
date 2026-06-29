package cmd

import (
	"os"
	"path/filepath"

	"github.com/jonamat/hetzner-rescaler/internal/store"
)

func openStore() (*store.Store, error) {
	path := os.Getenv("RESCALER_DB_PATH")
	if path == "" {
		path = filepath.Join(userHomeDir(), ".hetzner-rescaler", "db.sqlite")
	}
	return store.Open(path)
}

func userHomeDir() string {
	home, _ := os.UserHomeDir()
	if home == "" {
		return "."
	}
	return home
}
