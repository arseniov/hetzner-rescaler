package crypto

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
)

// Keyring holds the 32-byte AES-256 key used to seal Hetzner Cloud tokens
// at rest. A zero value is invalid; use NewKeyring (tests) or LoadKeyring
// (production) to obtain one.
type Keyring struct {
	key []byte
}

// NewKeyring returns a Keyring seeded with a fresh random 32-byte key.
// Intended for tests; production code should use LoadKeyring so the key
// is persistent across restarts.
func NewKeyring() (*Keyring, error) {
	k := make([]byte, 32)
	if _, err := rand.Read(k); err != nil {
		return nil, fmt.Errorf("crypto: read random key: %w", err)
	}
	return &Keyring{key: k}, nil
}

// NewKeyringFromBytes wraps an externally-sourced 32-byte key in a Keyring.
// Used by the cmd layer, which loads the key from RESCALER_TOKEN_ENCRYPTION_KEY
// or a file (raw bytes — not base64/hex). Returns ErrInvalidKey if the key
// is not exactly 32 bytes.
func NewKeyringFromBytes(raw []byte) (*Keyring, error) {
	if len(raw) != 32 {
		return nil, ErrInvalidKey
	}
	c := make([]byte, 32)
	copy(c, raw)
	return &Keyring{key: c}, nil
}

// LoadKeyring returns a Keyring from the RESCALER_ENCRYPTION_KEY
// environment variable. The variable can be base64 (preferred) or hex
// encoded; it must decode to exactly 32 bytes. If unset, a fresh random
// key is generated (test/no-config behaviour).
func LoadKeyring() (*Keyring, error) {
	raw := os.Getenv("RESCALER_ENCRYPTION_KEY")
	if raw == "" {
		return NewKeyring()
	}
	decoded, err := decodeKey(raw)
	if err != nil {
		return nil, err
	}
	if len(decoded) != 32 {
		return nil, ErrInvalidKey
	}
	return &Keyring{key: decoded}, nil
}

// Seal encrypts plaintext and returns (ciphertext, nonce). The nonce must
// be persisted alongside the ciphertext to allow decryption later.
func (k *Keyring) Seal(plaintext []byte) (ciphertext, nonce []byte, err error) {
	if k == nil || len(k.key) != 32 {
		return nil, nil, ErrInvalidKey
	}
	return Encrypt(k.key, plaintext)
}

// Open decrypts ciphertext with the stored nonce. Provided for callers
// that need to round-trip a token through the keyring.
func (k *Keyring) Open(ciphertext, nonce []byte) ([]byte, error) {
	if k == nil || len(k.key) != 32 {
		return nil, ErrInvalidKey
	}
	return Decrypt(k.key, ciphertext, nonce)
}

func decodeKey(raw string) ([]byte, error) {
	if b, err := base64.StdEncoding.DecodeString(raw); err == nil {
		return b, nil
	}
	if b, err := hex.DecodeString(raw); err == nil {
		return b, nil
	}
	return nil, errors.New("crypto: RESCALER_ENCRYPTION_KEY must be base64 or hex")
}
