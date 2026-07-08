// Package rescaler is the low-level "shutdown -> change-type -> poweron" flow.
// It is mode-agnostic; the scheduler decides when to call it and with which
// target type. It depends only on the hetzner.API interface.
package rescaler

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hetznercloud/hcloud-go/hcloud"
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
// Returns hetzner.ErrUnavailable (wrapped) if the target is out of stock.
// The phaseHook is called with the next phase label at each phase boundary
// (shutting_down, changing_type, powering_on, done). Pass nil for no-op
// behavior — preserves the original 3-arg call signature for tests.
func Rescale(ctx context.Context, api hetzner.API, srv *hetzner.Server, targetType string) error {
	return RescaleWithHook(ctx, api, srv, targetType, nil)
}

// RescaleWithHook is Rescale plus an optional phase callback. See Rescale.
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
	emit("done")
	return nil
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