package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/jonamat/hetzner-rescaler/internal/store"
)

func TestRescale_RequiresConfirm(t *testing.T) {
	deps, _ := newTestDeps(t)
	h := NewRouter(deps)
	_, sid := seedServer(t, deps, "p1", "web-1")

	body := RescaleRequest{Direction: "up", Confirm: false}
	req := authedRequest(t, "POST", "/api/servers/"+itoa(sid)+"/rescale", body)
	rr := recorder(t, h, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("want 400 without confirm, got %d", rr.Code)
	}
}

func TestRescale_RejectsBadDirection(t *testing.T) {
	deps, _ := newTestDeps(t)
	h := NewRouter(deps)
	_, sid := seedServer(t, deps, "p1", "web-1")

	body := RescaleRequest{Direction: "sideways", Confirm: true}
	req := authedRequest(t, "POST", "/api/servers/"+itoa(sid)+"/rescale", body)
	rr := recorder(t, h, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d", rr.Code)
	}
}

func TestRescale_UpTriggersRescale(t *testing.T) {
	deps, _ := newTestDeps(t)

	// Override the rescaler executor so the test does not depend on the
	// real scheduler dispatch path (which is integration-tested elsewhere).
	rescaleCalled := make(chan struct{}, 1)
	deps.Rescaler = func(ctx context.Context, srv *store.Server, target string) error {
		rescaleCalled <- struct{}{}
		return nil
	}

	h := NewRouter(deps)
	_, sid := seedServer(t, deps, "p1", "web-1")
	body := RescaleRequest{Direction: "up", Confirm: true}
	req := authedRequest(t, "POST", "/api/servers/"+itoa(sid)+"/rescale", body)
	rr := recorder(t, h, req)
	if rr.Code != http.StatusAccepted {
		t.Fatalf("want 202, got %d (body=%q)", rr.Code, rr.Body.String())
	}
	select {
	case <-rescaleCalled:
	case <-time.After(2 * time.Second):
		t.Fatalf("rescaler not invoked")
	}
}

func TestPromote_RequiresAutoPromoteMode(t *testing.T) {
	deps, _ := newTestDeps(t)
	h := NewRouter(deps)
	_, sid := seedServer(t, deps, "p1", "web-1") // mode is "manual" from seedServer

	body := ConfirmRequest{Confirm: true}
	req := authedRequest(t, "POST", "/api/servers/"+itoa(sid)+"/promote", body)
	rr := recorder(t, h, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("want 400 (wrong mode), got %d", rr.Code)
	}
}

func TestPromote_RequiresConfirm(t *testing.T) {
	deps, _ := newTestDeps(t)
	h := NewRouter(deps)
	_, sid := seedServer(t, deps, "p1", "web-1")

	body := ConfirmRequest{Confirm: false}
	req := authedRequest(t, "POST", "/api/servers/"+itoa(sid)+"/promote", body)
	rr := recorder(t, h, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d", rr.Code)
	}
}

func TestPromote_SetsPromoteState(t *testing.T) {
	deps, _ := newTestDeps(t)
	h := NewRouter(deps)
	_, sid := seedServer(t, deps, "p1", "web-1")

	// Switch server to auto_promote mode.
	srv, _ := deps.Store.GetServer(sid)
	srv.Mode = "auto_promote"
	if err := deps.Store.UpdateServer(*srv); err != nil {
		t.Fatalf("UpdateServer: %v", err)
	}

	body := ConfirmRequest{Confirm: true}
	req := authedRequest(t, "POST", "/api/servers/"+itoa(sid)+"/promote", body)
	rr := recorder(t, h, req)
	if rr.Code != http.StatusAccepted {
		t.Fatalf("want 202, got %d", rr.Code)
	}

	updated, _ := deps.Store.GetServer(sid)
	if updated.PromoteState == nil || *updated.PromoteState != "promote_requested" {
		t.Fatalf("expected promote_state=promote_requested, got %v", updated.PromoteState)
	}
}

func TestDemote_RequiresAutoPromoteMode(t *testing.T) {
	deps, _ := newTestDeps(t)
	h := NewRouter(deps)
	_, sid := seedServer(t, deps, "p1", "web-1")
	body := ConfirmRequest{Confirm: true}
	req := authedRequest(t, "POST", "/api/servers/"+itoa(sid)+"/demote", body)
	rr := recorder(t, h, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("want 400 (wrong mode), got %d", rr.Code)
	}
}

func TestRescale_FailsWhenRescalerErrors(t *testing.T) {
	deps, _ := newTestDeps(t)
	deps.Rescaler = func(ctx context.Context, srv *store.Server, target string) error {
		return errors.New("simulated failure")
	}
	h := NewRouter(deps)
	_, sid := seedServer(t, deps, "p1", "web-1")
	body := RescaleRequest{Direction: "up", Confirm: true}
	req := authedRequest(t, "POST", "/api/servers/"+itoa(sid)+"/rescale", body)
	rr := recorder(t, h, req)
	if rr.Code == http.StatusAccepted {
		t.Fatalf("expected non-202 on rescaler error, got 202")
	}
	var er map[string]string
	_ = json.Unmarshal(rr.Body.Bytes(), &er)
	if er["error"] == "" {
		t.Fatalf("expected error field in body, got %q", rr.Body.String())
	}
}
