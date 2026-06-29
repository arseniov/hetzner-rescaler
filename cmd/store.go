package cmd

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
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

func loadOrGenerateKey() ([]byte, error) {
	if envKey := os.Getenv("RESCALER_TOKEN_ENCRYPTION_KEY"); envKey != "" {
		if len(envKey) != 64 { // hex-encoded 32 bytes
			return nil, fmt.Errorf("RESCALER_TOKEN_ENCRYPTION_KEY must be 64 hex chars (32 bytes); got %d", len(envKey))
		}
		key, err := hex.DecodeString(envKey)
		if err != nil {
			return nil, fmt.Errorf("RESCALER_TOKEN_ENCRYPTION_KEY is not valid hex: %w", err)
		}
		return key, nil
	}
	dir := filepath.Join(userHomeDir(), ".hetzner-rescaler")
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, err
	}
	keyPath := filepath.Join(dir, "key")
	if data, err := os.ReadFile(keyPath); err == nil {
		if len(data) != 32 {
			return nil, fmt.Errorf("key file %s has wrong size: %d", keyPath, len(data))
		}
		return data, nil
	}
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return nil, err
	}
	if err := os.WriteFile(keyPath, key, 0600); err != nil {
		return nil, err
	}
	return key, nil
}
