package api

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestProjectResponse_JSONShape(t *testing.T) {
	p := ProjectResponse{
		ID:           7,
		Name:         "prod",
		HasToken:     true,
		LastError:    "",
		CreatedAt:    time.Date(2026, 6, 29, 12, 0, 0, 0, time.UTC),
		UpdatedAt:    time.Date(2026, 6, 29, 12, 0, 0, 0, time.UTC),
	}
	out, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	for _, want := range []string{
		`"id":7`,
		`"name":"prod"`,
		`"has_token":true`,
		`"last_error":""`,
		`"created_at":"2026-06-29T12:00:00Z"`,
	} {
		if !strings.Contains(string(out), want) {
			t.Errorf("JSON missing %q; got %s", want, out)
		}
	}
}

func TestServerResponse_JSONShape(t *testing.T) {
	s := ServerResponse{
		ID:             42,
		ProjectID:      7,
		HCloudServerID: 12345,
		Name:           "web-1",
		Label:          "primary web",
		BaseServerType: "cpx11",
		TopServerType:  "cpx31",
		FallbackChain:  []string{"cpx31", "cpx21", "cpx11"},
		Mode:           "scheduled",
		PromoteState:   nil,
		Timezone:       "Europe/Rome",
	}
	out, err := json.Marshal(s)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	for _, want := range []string{
		`"id":42`,
		`"project_id":7`,
		`"hcloud_server_id":12345`,
		`"name":"web-1"`,
		`"base_server_type":"cpx11"`,
		`"top_server_type":"cpx31"`,
		`"fallback_chain":["cpx31","cpx21","cpx11"]`,
		`"mode":"scheduled"`,
		`"timezone":"Europe/Rome"`,
	} {
		if !strings.Contains(string(out), want) {
			t.Errorf("JSON missing %q; got %s", want, out)
		}
	}
}

func TestRescaleRequest_RequiresConfirm(t *testing.T) {
	var r RescaleRequest
	if err := json.Unmarshal([]byte(`{"direction":"up"}`), &r); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if r.Confirm {
		t.Fatalf("expected confirm=false when omitted")
	}
}

func TestConfirmRequest_DefaultFalse(t *testing.T) {
	var r ConfirmRequest
	if r.Confirm {
		t.Fatalf("zero value should be false")
	}
}