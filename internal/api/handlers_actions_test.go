package api

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/jonamat/hetzner-rescaler/internal/hetzner"
	"github.com/jonamat/hetzner-rescaler/internal/rescaler"
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

func TestRescale_UpReturns202WithPendingEventID(t *testing.T) {
	deps, _ := newTestDeps(t)

	mgr := rescaler.NewManager(deps.Store)
	if err := mgr.Start(context.Background()); err != nil {
		t.Fatalf("manager Start: %v", err)
	}
	// Empty fakeHetzner — GetServer returns nil so Submit records FromType="".
	// The goroutine will panic asynchronously; we only assert on the
	// synchronous Submit return.
	mgr.SetAPIResolver(func(ctx context.Context, projectID int64) (hetzner.API, error) {
		return &fakeHetzner{}, nil
	})
	deps.Manager = mgr

	h := NewRouter(deps)
	_, sid := seedServer(t, deps, "p1", "web-1")
	body := RescaleRequest{Direction: "up", Confirm: true}
	req := authedRequest(t, "POST", "/api/servers/"+itoa(sid)+"/rescale", body)
	rr := recorder(t, h, req)
	if rr.Code != http.StatusAccepted {
		t.Fatalf("want 202, got %d (body=%q)", rr.Code, rr.Body.String())
	}
	var resp map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode: %v (body=%q)", err, rr.Body.String())
	}
	if resp["status"] != "rescale initiated" {
		t.Fatalf("status = %v, want 'rescale initiated'", resp["status"])
	}
	pid, ok := resp["pending_event_id"].(float64)
	if !ok || int64(pid) == 0 {
		t.Fatalf("pending_event_id missing or zero in %v", resp)
	}

	// Drain goroutine (it will panic on the fakeHetzner; recover via Shutdown).
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_ = mgr.Shutdown(shutdownCtx)
}

func TestRescale_UpTriggersRescale(t *testing.T) {
	deps, _ := newTestDeps(t)

	mgr := rescaler.NewManager(deps.Store)
	if err := mgr.Start(context.Background()); err != nil {
		t.Fatalf("manager Start: %v", err)
	}
	// Stub resolver returning a fakeHetzner whose GetServer returns nil
	// (empty fakeHetzner). The goroutine will eventually panic, but the
	// assertion below only cares about the synchronous onSubmit hook.
	mgr.SetAPIResolver(func(ctx context.Context, projectID int64) (hetzner.API, error) {
		return &fakeHetzner{}, nil
	})
	// SetOnSubmit fires BEFORE the busy check, so the channel-send is
	// visible to the test even before Submit's spawned goroutine
	// progresses. This is what makes the assertion deterministic.
	called := make(chan struct{}, 1)
	mgr.SetOnSubmit(func() { called <- struct{}{} })

	deps.Manager = mgr

	h := NewRouter(deps)
	_, sid := seedServer(t, deps, "p1", "web-1")
	body := RescaleRequest{Direction: "up", Confirm: true}
	req := authedRequest(t, "POST", "/api/servers/"+itoa(sid)+"/rescale", body)
	rr := recorder(t, h, req)
	if rr.Code != http.StatusAccepted {
		t.Fatalf("want 202, got %d (body=%q)", rr.Code, rr.Body.String())
	}
	select {
	case <-called:
	case <-time.After(2 * time.Second):
		t.Fatalf("manager onSubmit hook not invoked")
	}

	// Drain goroutine.
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_ = mgr.Shutdown(shutdownCtx)
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
