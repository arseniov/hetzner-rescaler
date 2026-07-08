package cmd

import (
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jonamat/hetzner-rescaler/internal/crypto"
	"github.com/jonamat/hetzner-rescaler/internal/store"
)

func openStore() (*store.Store, error) {
	path := os.Getenv("RESCALER_DB_PATH")
	if path == "" {
		path = filepath.Join(userHomeDir(), ".hetzner-rescaler", "db.sqlite")
	}
	return store.Open(path)
}

// openStoreAndKeyring opens the SQLite store at the configured path and
// builds (or loads) an AES-GCM keyring keyed for token encryption. Used
// by `serve` which needs both at startup.
func openStoreAndKeyring() (*store.Store, *crypto.Keyring, error) {
	raw, err := loadOrGenerateKey()
	if err != nil {
		return nil, nil, fmt.Errorf("load key: %w", err)
	}
	key, err := crypto.NewKeyringFromBytes(raw)
	if err != nil {
		return nil, nil, fmt.Errorf("init keyring: %w", err)
	}
	s, err := openStore()
	if err != nil {
		return nil, nil, err
	}
	return s, key, nil
}

func userHomeDir() string {
	home, _ := os.UserHomeDir()
	if home == "" {
		return "."
	}
	return home
}

// loadOrGenerateKey reads the AES-GCM encryption key used to seal Hetzner
// tokens at rest in the SQLite store.
//
// RESCALER_TOKEN_ENCRYPTION_KEY is REQUIRED. We deliberately do not
// auto-generate and persist a key (the previous behaviour) because that
// path is silent: an operator who rebuilds the Docker image without
// setting the env var gets a fresh key, then every previously-stored
// token fails to decrypt with `message authentication failed` (a GCM
// authentication tag mismatch). Failing loudly at startup is the only
// way to surface that mistake before any HTTP traffic hits the API.
//
// To generate a key:
//   openssl rand -hex 32
// then store the 64-char hex string in the env var.
func loadOrGenerateKey() ([]byte, error) {
	envKey := os.Getenv("RESCALER_TOKEN_ENCRYPTION_KEY")
	if envKey == "" {
		return nil, fmt.Errorf(
			"RESCALER_TOKEN_ENCRYPTION_KEY is required; generate one with `openssl rand -hex 32` " +
				"and set it in the environment before running any command that touches tokens " +
				"(serve, config, start, try, migrate). Without a stable key, tokens already " +
				"encrypted into the SQLite store become unreadable on the next process restart.",
		)
	}
	if len(envKey) != 64 { // hex-encoded 32 bytes
		return nil, fmt.Errorf("RESCALER_TOKEN_ENCRYPTION_KEY must be 64 hex chars (32 bytes); got %d", len(envKey))
	}
	key, err := hex.DecodeString(envKey)
	if err != nil {
		return nil, fmt.Errorf("RESCALER_TOKEN_ENCRYPTION_KEY is not valid hex: %w", err)
	}
	return key, nil
}
