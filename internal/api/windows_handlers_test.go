package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/jonamat/hetzner-rescaler/internal/store"
)

func TestListWindows_ReturnsArrayForServer(t *testing.T) {
	deps, _ := newTestDeps(t)
	h := NewRouter(deps)
	_, sid := seedServer(t, deps, "p1", "web-1")

	// Empty list initially.
	req := authedRequest(t, "GET", "/api/servers/"+itoa(sid)+"/windows", nil)
	rr := recorder(t, h, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", rr.Code)
	}
	body := strings.TrimSpace(rr.Body.String())
	if body != "[]" {
		t.Fatalf("want [] for empty list, got %q", body)
	}

	// Add one window via POST.
	post := CreateWindowRequest{
		Label: "weekday-peak", DaysOfWeek: 0b00111110, StartTime: "09:00", StopTime: "18:00",
		TargetType: "cpx31", Enabled: true,
	}
	req = authedRequest(t, "POST", "/api/servers/"+itoa(sid)+"/windows", post)
	rr = recorder(t, h, req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("want 201, got %d (body=%q)", rr.Code, rr.Body.String())
	}

	// GET again — should now have one entry.
	req = authedRequest(t, "GET", "/api/servers/"+itoa(sid)+"/windows", nil)
	rr = recorder(t, h, req)
	var list []WindowResponse
	_ = json.Unmarshal(rr.Body.Bytes(), &list)
	if len(list) != 1 || list[0].Label != "weekday-peak" {
		t.Fatalf("unexpected: %+v", list)
	}
}

func TestCreateWindow_ValidatesTimes(t *testing.T) {
	deps, _ := newTestDeps(t)
	h := NewRouter(deps)
	_, sid := seedServer(t, deps, "p1", "web-1")

	cases := []CreateWindowRequest{
		{Label: "x", DaysOfWeek: 0, StartTime: "bad", StopTime: "10:00", TargetType: "cpx31"},
		{Label: "x", DaysOfWeek: 0, StartTime: "10:00", StopTime: "10:00", TargetType: "cpx31"},
		{Label: "", DaysOfWeek: 0, StartTime: "10:00", StopTime: "12:00", TargetType: "cpx31"},
		{Label: "x", DaysOfWeek: 0, StartTime: "10:00", StopTime: "11:00", TargetType: ""},
	}
	for i, body := range cases {
		req := authedRequest(t, "POST", "/api/servers/"+itoa(sid)+"/windows", body)
		rr := recorder(t, h, req)
		if rr.Code == http.StatusCreated {
			t.Fatalf("case %d expected non-201, got 201", i)
		}
	}
}

func TestUpdateWindow_AppliesChanges(t *testing.T) {
	deps, _ := newTestDeps(t)
	h := NewRouter(deps)
	_, sid := seedServer(t, deps, "p1", "web-1")
	w, err := deps.Store.CreateWindow(sid, store.Window{
		Label: "orig", DaysOfWeek: 0b00111110, StartTime: "09:00",
		StopTime: "18:00", TargetType: "cpx31", Enabled: true,
	})
	if err != nil {
		t.Fatalf("CreateWindow: %v", err)
	}

	body := UpdateWindowRequest{
		Label: "renamed", DaysOfWeek: 0b01000001, StartTime: "10:00",
		StopTime: "20:00", TargetType: "cpx21", Enabled: false,
	}
	req := authedRequest(t, "PUT", "/api/windows/"+itoa(w.ID), body)
	rr := recorder(t, h, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("want 200, got %d (body=%q)", rr.Code, rr.Body.String())
	}
	var got WindowResponse
	_ = json.Unmarshal(rr.Body.Bytes(), &got)
	if got.Label != "renamed" || got.DaysOfWeek != 0b01000001 || got.TargetType != "cpx21" {
		t.Fatalf("update did not apply: %+v", got)
	}
}

func TestDeleteWindow_RemovesRow(t *testing.T) {
	deps, _ := newTestDeps(t)
	h := NewRouter(deps)
	_, sid := seedServer(t, deps, "p1", "web-1")
	w, _ := deps.Store.CreateWindow(sid, store.Window{
		Label: "del", DaysOfWeek: 0b00111110, StartTime: "09:00", StopTime: "18:00", TargetType: "cpx31",
	})
	req := authedRequest(t, "DELETE", "/api/windows/"+itoa(w.ID), nil)
	rr := recorder(t, h, req)
	if rr.Code != http.StatusNoContent {
		t.Fatalf("want 204, got %d", rr.Code)
	}
}
