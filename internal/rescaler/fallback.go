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
func RescaleWithFallback(ctx context.Context, api hetzner.API, srv *hetzner.Server, targetType string, chain []string) (string, error) {
	return RescaleWithFallbackWithHook(ctx, api, srv, targetType, chain, nil)
}

// RescaleWithFallbackWithHook is RescaleWithFallback with a phase hook
// passed down to every Rescale invocation. Same return contract.
func RescaleWithFallbackWithHook(ctx context.Context, api hetzner.API, srv *hetzner.Server, targetType string, chain []string, phaseHook func(string)) (string, error) {
	for _, t := range chain {
		err := RescaleWithHook(ctx, api, srv, t, phaseHook)
		if err == nil {
			return t, nil
		}
		if hetzner.IsUnavailable(err) {
			continue
		}
		return "", fmt.Errorf("rescaler: chain halted at %q: %w", t, err)
	}
	return "", ErrAllUnavailable
}