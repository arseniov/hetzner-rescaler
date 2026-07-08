package rescaler

import (
	"context"
	"errors"
	"fmt"

	"github.com/jonamat/hetzner-rescaler/internal/hetzner"
)

// ErrAllUnavailable is returned by RescaleWithFallback when every type in
// the chain is unavailable. The caller should treat this as a transient
// failure (e.g. surface as a failed event) and retry on the next window.
var ErrAllUnavailable = errors.New("rescaler: all fallback targets unavailable")

// RescaleWithFallback tries targetType first; if it returns an unavailable
// error, it walks the chain. Returns the type that actually succeeded, or
// ErrAllUnavailable if the chain was exhausted, or a non-unavailable error
// from the first failing attempt (subsequent chain entries are NOT tried
// for non-unavailable errors).
//
// Contract: when RescaleWithFallback returns, the server is ALWAYS
// running — if a chain entry succeeded, its success path already
// powered the server on; if all entries failed (or the chain halted
// on a non-unavailable error), the safety-net poweron below
// recovers the server. Operators should never have to chase a
// production server that was left off after a "failed rescale".
func RescaleWithFallback(ctx context.Context, api hetzner.API, srv *hetzner.Server, targetType string, chain []string) (string, error) {
	return RescaleWithFallbackWithHook(ctx, api, srv, targetType, chain, nil)
}

// RescaleWithFallbackWithHook is RescaleWithFallback with a phase hook
// passed down to every Rescale invocation. Same return contract.
//
// The chain is walked via RescaleWithHook, which does NOT auto-poweron
// on failure (the wrapper's safety net handles that). This split
// keeps the chain from paying a shutdown/poweron cycle per iteration
// when every entry is unavailable — only the final attempt gets the
// recovery poweron.
func RescaleWithFallbackWithHook(ctx context.Context, api hetzner.API, srv *hetzner.Server, targetType string, chain []string, phaseHook func(string)) (string, error) {
	for _, t := range chain {
		err := RescaleWithHook(ctx, api, srv, t, phaseHook)
		if err == nil {
			return t, nil
		}
		if hetzner.IsUnavailable(err) {
			continue
		}
		// Non-unavailable error: halt the chain. The server may be in
		// any state (e.g. off if shutdown succeeded but change_type
		// hit something unexpected), so the safety net is essential.
		ensurePoweredOn(ctx, api, srv)
		return "", fmt.Errorf("rescaler: chain halted at %q: %w", t, err)
	}
	// Chain exhausted on unavailable errors only. Server is off
	// (we shut it down in the first iteration, never powered it back).
	// Recover.
	ensurePoweredOn(ctx, api, srv)
	return "", ErrAllUnavailable
}