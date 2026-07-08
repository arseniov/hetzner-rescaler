package cmd

import (
	"encoding/hex"
	"strings"
	"testing"
)

func TestLoadOrGenerateKey_RequiresEnvVar(t *testing.T) {
	// Unset the env var so we can prove the loader refuses to fall
	// back to a generated key. (The t.Setenv("...", "") idiom is
	// documented to clear the variable for the test's duration.)
	t.Setenv("RESCALER_TOKEN_ENCRYPTION_KEY", "")

	_, err := loadOrGenerateKey()
	if err == nil {
		t.Fatal("loadOrGenerateKey() with empty env var: want error, got nil")
	}
	// Error must mention the env var name so the operator knows
	// exactly what to set.
	if !strings.Contains(err.Error(), "RESCALER_TOKEN_ENCRYPTION_KEY") {
		t.Errorf("error must mention RESCALER_TOKEN_ENCRYPTION_KEY; got: %v", err)
	}
	// Error must include the generation recipe so the operator
	// can fix it without reading the docs.
	if !strings.Contains(err.Error(), "openssl rand -hex 32") {
		t.Errorf("error must include the generation command; got: %v", err)
	}
}

func TestLoadOrGenerateKey_AcceptsValidHex(t *testing.T) {
	key := make([]byte, 32)
	for i := range key {
		key[i] = byte(i)
	}
	t.Setenv("RESCALER_TOKEN_ENCRYPTION_KEY", hex.EncodeToString(key))

	got, err := loadOrGenerateKey()
	if err != nil {
		t.Fatalf("loadOrGenerateKey() with valid env var: %v", err)
	}
	if len(got) != 32 {
		t.Errorf("len(got) = %d, want 32", len(got))
	}
	for i, b := range got {
		if b != key[i] {
			t.Errorf("key byte %d = %d, want %d", i, b, key[i])
		}
	}
}

func TestLoadOrGenerateKey_RejectsBadHexLength(t *testing.T) {
	t.Setenv("RESCALER_TOKEN_ENCRYPTION_KEY", "abc123")
	_, err := loadOrGenerateKey()
	if err == nil {
		t.Fatal("loadOrGenerateKey() with 6-char hex: want error, got nil")
	}
	if !strings.Contains(err.Error(), "64 hex chars") {
		t.Errorf("error must mention the expected length; got: %v", err)
	}
}

func TestLoadOrGenerateKey_RejectsInvalidHex(t *testing.T) {
	// 64 chars but contains non-hex characters
	bad := "zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz"
	t.Setenv("RESCALER_TOKEN_ENCRYPTION_KEY", bad)
	_, err := loadOrGenerateKey()
	if err == nil {
		t.Fatal("loadOrGenerateKey() with non-hex chars: want error, got nil")
	}
	if !strings.Contains(err.Error(), "not valid hex") {
		t.Errorf("error must mention invalid hex; got: %v", err)
	}
}