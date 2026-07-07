package rescaler

import (
	"testing"

	"github.com/jonamat/hetzner-rescaler/internal/store"
)

func TestNewManager_NotNil(t *testing.T) {
	s, err := store.OpenTemp()
	if err != nil {
		t.Fatalf("OpenTemp: %v", err)
	}
	defer s.Close()
	m := NewManager(s)
	if m == nil {
		t.Fatal("NewManager returned nil")
	}
	if m.jobs == nil {
		t.Fatal("jobs map not initialised")
	}
}