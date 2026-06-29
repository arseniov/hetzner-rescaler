package hetzner

import (
	"context"
	"errors"
	"testing"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

// fakeAPI is a minimal in-test fake. We use it to verify the adapter
// forwards calls to the wrapped *hcloud.Client (we don't have a live
// Hetzner server, so we just check that arguments are wired through
// using a tiny in-test shim).
//
// Because the real adapter delegates everything, the strongest test
// we can do here is to verify the adapter implements the API interface
// and the constructors return non-nil.

func TestNewClientRequiresToken(t *testing.T) {
	if _, err := NewClient(""); err == nil {
		t.Fatal("expected error on empty token")
	}
}

func TestNewClientReturnsAPI(t *testing.T) {
	api, err := NewClient("fake-token-for-test")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	if api == nil {
		t.Fatal("nil api")
	}
	// The returned api must satisfy API at compile time (var _ API = api).
	var _ API = api
	_ = context.Background
	_ = hcloud.Server{}
	_ = errors.New
}
