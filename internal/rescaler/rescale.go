// Package rescaler is the low-level "shutdown -> change-type -> poweron" flow.
// It is mode-agnostic; the scheduler decides when to call it and with which
// target type. It depends only on the hetzner.API interface.
package rescaler

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/jonamat/hetzner-rescaler/internal/hetzner"
)

// pollInterval and provisionerSleep match the original implementation's
// behavior. Tune in one place.
const (
	pollInterval     = 5 * time.Second
	provisionerSleep = 30 * time.Second
)

// Rescale changes the given server to targetType. The flow:
//  1. if server is already at targetType, no-op
//  2. shutdown the server (if running)
//  3. wait for the Hetzner provisioner
//  4. change type
//  5. poweron
//
// Contract: when Rescale returns, the server is ALWAYS running —
// regardless of whether the rescale succeeded, the target was
// unavailable, or the change_type action failed. If anything in the
// middle of the flow goes wrong, the safety-net poweron at the end
// recovers the server. An off server is unreachable, so a failed
// rescale that silently left the server off was indistinguishable
// from a production outage — exactly the failure mode the safety
// net is built to prevent.
//
// Returns hetzner.ErrUnavailable (wrapped) if the target is out of
// stock. The phaseHook is called with the next phase label at each
// phase boundary (shutting_down, changing_type, powering_on, done).
func Rescale(ctx context.Context, api hetzner.API, srv *hetzner.Server, targetType string) error {
	err := RescaleWithHook(ctx, api, srv, targetType, nil)
	// Safety net: if the inner flow returned an error (target
	// unavailable, change_type action failed, etc.), the server is
	// off — recover it. On success the inner flow already powered on,
	// so this is a no-op.
	if err != nil {
		ensurePoweredOn(ctx, api, srv)
	}
	return err
}

// RescaleWithHook is the inner flow: shutdown → change_type → poweron
// on the success path only. On any error path it returns WITHOUT
// powering the server back on — the caller (Rescale or
// RescaleWithFallbackWithHook) is responsible for the safety-net
// poweron. Splitting concerns this way keeps the fallback chain from
// paying redundant shutdown/poweron cycles between iterations.
//
// Use Rescale (or RescaleWithFallbackWithHook) from production code;
// RescaleWithHook is exported so tests and other inner flows can
// drive a single attempt without the wrapper's safety net.
func RescaleWithHook(ctx context.Context, api hetzner.API, srv *hetzner.Server, targetType string, phaseHook func(string)) error {
	emit := func(phase string) {
		if phaseHook != nil {
			phaseHook(phase)
		}
	}

	if srv.ServerType.Name == targetType {
		return nil
	}

	// We only sleep the provisioner wait when transitioning from running →
	// change_type. Off → change_type intentionally omits the sleep because no
	// shutdown just happened for the provisioner to recover from.
	if srv.Status == hcloud.ServerStatusRunning {
		emit("shutting_down")
		act, err := api.ShutdownServer(ctx, srv)
		if err != nil {
			return fmt.Errorf("rescaler: shutdown: %w", err)
		}
		if err := waitAction(ctx, api, act); err != nil {
			return fmt.Errorf("rescaler: wait shutdown: %w", err)
		}
		// Reflect the new state on the in-memory server struct so the
		// fallback chain (and any caller polling this pointer) sees a
		// consistent picture. Without this, a fallback loop would try
		// to shut down the same already-off server again on the next
		// iteration — paying another 30s provisioner wait and risking
		// a "server already off" error from Hetzner.
		srv.Status = hcloud.ServerStatusOff
		// Wait for the Hetzner provisioner to be ready to accept a new type
		time.Sleep(provisionerSleep)
	}

	target, err := api.GetServerType(ctx, targetType)
	if err != nil {
		return fmt.Errorf("rescaler: get target type: %w", err)
	}
	emit("changing_type")
	act, err := api.ChangeServerType(ctx, srv, target)
	if err != nil {
		return err // may be wrapped hetzner.ErrUnavailable; caller classifies
	}
	if err := waitAction(ctx, api, act); err != nil {
		return fmt.Errorf("rescaler: wait change_type: %w", err)
	}

	// Reflect the new type on the in-memory server struct. (We don't refetch
	// from Hetzner here; the caller can do that if it needs the canonical
	// state.)
	srv.ServerType = target

	emit("powering_on")
	act, err = api.PowerOnServer(ctx, srv)
	if err != nil {
		return fmt.Errorf("rescaler: poweron: %w", err)
	}
	if err := waitAction(ctx, api, act); err != nil {
		return fmt.Errorf("rescaler: wait poweron: %w", err)
	}
	srv.Status = hcloud.ServerStatusRunning
	emit("done")
	return nil
}

// ensurePoweredOn is the safety net invoked by Rescale and
// RescaleWithFallbackWithHook when the inner flow left the server
// off. It powers the server back on and updates the in-memory
// status. Errors are swallowed: at this point we've already decided
// to return a primary error to the caller, and a poweron failure
// shouldn't shadow the original cause.
//
// Uses context.Background so the recovery isn't tied to a possibly-
// cancelled caller's context (e.g. an HTTP client that closed).
func ensurePoweredOn(ctx context.Context, api hetzner.API, srv *hetzner.Server) {
	_ = ctx // currently unused — we deliberately use Background below
	if srv.Status == hcloud.ServerStatusRunning {
		return
	}
	act, err := api.PowerOnServer(context.Background(), srv)
	if err != nil {
		return
	}
	if err := waitAction(context.Background(), api, act); err != nil {
		return
	}
	srv.Status = hcloud.ServerStatusRunning
}

func waitAction(ctx context.Context, api hetzner.API, act *hetzner.Action) error {
	for {
		got, err := api.GetAction(ctx, act.ID)
		if err != nil {
			return err
		}
		switch got.Status {
		case hcloud.ActionStatusError:
			if aerr := got.Error(); aerr != nil {
				return aerr
			}
			return errors.New("rescaler: action failed without error message")
		case hcloud.ActionStatusSuccess:
			return nil
		default:
			time.Sleep(pollInterval)
		}
	}
}