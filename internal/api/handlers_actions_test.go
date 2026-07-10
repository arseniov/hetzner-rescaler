package api

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
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

func TestRescale_ConcurrentReturns409WithPendingID(t *testing.T) {
	deps, _ := newTestDeps(t)
	_, sid := seedServer(t, deps, "p1", "web-1")

	deps.Manager = rescaler.NewManager(deps.Store)
	_ = deps.Manager.Start(context.Background())
	// Use a blocking API so the first rescale is still in-flight when the
	// second arrives.
	deps.Manager.SetAPIResolver(func(ctx context.Context, projectID int64) (hetzner.API, error) {
		return &blockingTestAPI{}, nil
	})

	h := NewRouter(deps)
	body := RescaleRequest{Direction: "up", Confirm: true}

	// First request — should return 202 and leave the goroutine in-flight.
	firstReq := authedRequest(t, "POST", "/api/servers/"+itoa(sid)+"/rescale", body)
	firstRR := recorder(t, h, firstReq)
	if firstRR.Code != http.StatusAccepted {
		t.Fatalf("first: want 202, got %d", firstRR.Code)
	}

	// Give the goroutine a moment to enter the rescale.
	time.Sleep(50 * time.Millisecond)

	// Second request — should return 409.
	secondReq := authedRequest(t, "POST", "/api/servers/"+itoa(sid)+"/rescale", body)
	secondRR := recorder(t, h, secondReq)
	if secondRR.Code != http.StatusConflict {
		t.Fatalf("second: want 409, got %d (body=%q)", secondRR.Code, secondRR.Body.String())
	}
	var resp map[string]any
	if err := json.Unmarshal(secondRR.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp["error"] != "rescale already in progress" {
		t.Fatalf("error = %v, want 'rescale already in progress'", resp["error"])
	}
	if id, ok := resp["pending_event_id"].(float64); !ok || id == 0 {
		t.Fatalf("pending_event_id missing: %v", resp["pending_event_id"])
	}

	// Clean up: cancel the in-flight goroutine so the test exits cleanly.
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_ = deps.Manager.Shutdown(ctx)
	})
}

type blockingTestAPI struct{ hetzner.API }

func (b *blockingTestAPI) GetServer(ctx context.Context, id int64) (*hetzner.Server, error) {
	// Status: ServerStatusRunning so RescaleWithHook hits the shutdown branch
	// (calling our ShutdownServer), then waits in waitAction's blocking GetAction.
	return &hetzner.Server{ID: id, Status: hcloud.ServerStatusRunning, ServerType: &hetzner.ServerType{Name: "cpx11"}}, nil
}

func (b *blockingTestAPI) ShutdownServer(ctx context.Context, srv *hetzner.Server) (*hetzner.Action, error) {
	return &hetzner.Action{ID: 1, Status: hcloud.ActionStatusRunning, Command: "shutdown"}, nil
}

func (b *blockingTestAPI) GetAction(ctx context.Context, id int64) (*hetzner.Action, error) {
	<-ctx.Done()
	return nil, ctx.Err()
}
