package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/jonamat/hetzner-rescaler/internal/store"
)

func TestStatusPrintsProjectsAndServers(t *testing.T) {
	st, err := store.OpenTemp()
	if err != nil {
		t.Fatalf("OpenTemp: %v", err)
	}
	defer st.Close()

	p, _ := st.CreateProject("alpha", []byte("tok"), []byte("nonce12byts"))
	_, _ = st.CreateServer(p.ID, store.Server{
		HCloudServerID: 1, Name: "web", BaseServerType: "cpx11", TopServerType: "cpx21",
		FallbackChain: []string{"cpx21", "cpx11"}, Mode: "scheduled", Timezone: "UTC",
	})

	var buf bytes.Buffer
	if err := printStatus(st, &buf); err != nil {
		t.Fatalf("printStatus: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "alpha") || !strings.Contains(out, "web") {
		t.Fatalf("output missing expected text: %q", out)
	}
}
