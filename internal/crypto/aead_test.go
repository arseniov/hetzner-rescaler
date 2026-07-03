package crypto

import (
	"bytes"
	"errors"
	"testing"
)

func TestEncryptDecryptRoundTrip(t *testing.T) {
	key := bytes.Repeat([]byte{0x42}, 32) // 32 bytes -> AES-256
	plaintext := []byte("hcloud_token_abc123")

	ciphertext, nonce, err := Encrypt(key, plaintext)
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}
	if bytes.Equal(ciphertext, plaintext) {
		t.Fatalf("ciphertext equals plaintext: encryption did not happen")
	}
	if len(nonce) != 12 {
		t.Fatalf("nonce length = %d, want 12", len(nonce))
	}

	got, err := Decrypt(key, ciphertext, nonce)
	if err != nil {
		t.Fatalf("Decrypt: %v", err)
	}
	if !bytes.Equal(got, plaintext) {
		t.Fatalf("decrypted = %q, want %q", got, plaintext)
	}
}

func TestEncryptProducesDifferentCiphertextsForSamePlaintext(t *testing.T) {
	key := bytes.Repeat([]byte{0x42}, 32)
	plaintext := []byte("same input")

	c1, _, _ := Encrypt(key, plaintext)
	c2, _, _ := Encrypt(key, plaintext)

	if bytes.Equal(c1, c2) {
		t.Fatalf("two encryptions of same plaintext produced identical ciphertexts (nonce reuse?)")
	}
}

func TestDecryptWithWrongKeyFails(t *testing.T) {
	key := bytes.Repeat([]byte{0x42}, 32)
	wrong := bytes.Repeat([]byte{0x43}, 32)
	plaintext := []byte("secret")

	ct, nonce, err := Encrypt(key, plaintext)
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}

	_, err = Decrypt(wrong, ct, nonce)
	if err == nil {
		t.Fatalf("Decrypt with wrong key succeeded, want error")
	}
}

func TestEncryptKeyLengthValidation(t *testing.T) {
	bad := bytes.Repeat([]byte{0x42}, 16) // 16 bytes is fine for AES-128, but we standardize on 32
	_, _, err := Encrypt(bad, []byte("x"))
	if err == nil {
		t.Fatalf("Encrypt with non-32-byte key succeeded, want error")
	}
	if !errors.Is(err, ErrInvalidKey) {
		t.Fatalf("err = %v, want ErrInvalidKey", err)
	}
}
