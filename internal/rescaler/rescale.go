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
	pollInterval      = 5 * time.Second
	provisionerSleep  = 30 * time.Second
)

// Rescale changes the given server to targetType. The flow:
//  1. if server is already at targetType, no-op
//  2. shutdown the server (if running)
//  3. wait for the Hetzner provisioner
//  4. change type
//  5. poweron
//
// Returns hetzner.ErrUnavailable (wrapped) if the target is out of stock.
func Rescale(ctx context.Context, api hetzner.API, srv *hetzner.Server, targetType string) error {
	if srv.ServerType.Name == targetType {
		return nil
	}

	if srv.Status == hcloud.ServerStatusRunning {
		act, err := api.ShutdownServer(ctx, srv)
		if err != nil {
			return fmt.Errorf("rescaler: shutdown: %w", err)
		}
		if err := waitAction(ctx, api, act); err != nil {
			return fmt.Errorf("rescaler: wait shutdown: %w", err)
		}
		// Wait for the Hetzner provisioner to be ready to accept a new type
		time.Sleep(provisionerSleep)
	}

	target, err := api.GetServerType(ctx, targetType)
	if err != nil {
		return fmt.Errorf("rescaler: get target type: %w", err)
	}
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

	act, err = api.PowerOnServer(ctx, srv)
	if err != nil {
		return fmt.Errorf("rescaler: poweron: %w", err)
	}
	if err := waitAction(ctx, api, act); err != nil {
		return fmt.Errorf("rescaler: wait poweron: %w", err)
	}
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
			return errors.New(got.ErrorMessage)
		case hcloud.ActionStatusSuccess:
			return nil
		default:
			time.Sleep(pollInterval)
		}
	}
}