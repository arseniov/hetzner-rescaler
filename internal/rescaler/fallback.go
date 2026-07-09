package rescaler

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jonamat/hetzner-rescaler/internal/hetzner"
	"github.com/jonamat/hetzner-rescaler/internal/store"
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
//
// This is the no-event-emission wrapper used by cmd/try.go. The scheduler
// and any future caller that wants rescale_skipped events in the store
// should call RescaleWithFallbackWithHook directly with a non-nil store
// and a valid local serverID.
func RescaleWithFallback(ctx context.Context, api hetzner.API, srv *hetzner.Server, targetType string, chain []string) (string, error) {
	return RescaleWithFallbackWithHook(ctx, api, srv, targetType, chain, nil, 0, "", nil)
}

// RescaleWithFallbackWithHook is RescaleWithFallback with a phase hook
// AND a per-entry availability pre-check.
//
// For each entry in chain, the loop first calls IsTypeAvailable:
//   - unavailable: emit a rescale_skipped event (if st != nil) and continue.
//   - check error: log a warning and fall through (fail open — see Risks).
//   - available:    run RescaleWithHook as today.
//
// The chain is walked via RescaleWithHook, which does NOT auto-poweron
// on failure (the wrapper's safety net handles that). This split
// keeps the chain from paying a shutdown/poweron cycle per iteration
// when every entry is unavailable — only the final attempt gets the
// recovery poweron.
//
// `st` may be nil; in that case no rescale_skipped events are emitted
// (the wrapper uses this to keep cmd/try.go as a no-event one-off).
// `triggeredBy` is recorded on the emitted events; pass "" when st is nil.
// `serverID` is the local store Server.ID (NOT the hcloud server ID) —
// the events table foreign-keys to servers.id, so we need the local row.
func RescaleWithFallbackWithHook(
	ctx context.Context,
	api hetzner.API,
	srv *hetzner.Server,
	targetType string,
	chain []string,
	st *store.Store,
	serverID int64,
	triggeredBy string,
	phaseHook func(string),
) (string, error) {
	for _, t := range chain {
		// IsTypeAvailable short-circuits to (false, nil) when the
		// server has no Datacenter/Location. That's a "we don't know"
		// signal, not a confirmed-unavailable — so we fail open
		// here too. Otherwise a freshly-observed server whose first
		// reconciliation hasn't completed would never rescale.
		if srv == nil || srv.Datacenter == nil || srv.Datacenter.Location == nil {
			err := RescaleWithHook(ctx, api, srv, t, phaseHook)
			if err == nil {
				return t, nil
			}
			if hetzner.IsUnavailable(err) {
				continue
			}
			ensurePoweredOn(ctx, api, srv)
			return "", fmt.Errorf("rescaler: chain halted at %q: %w", t, err)
		}
		if avail, err := IsTypeAvailable(ctx, api, srv, t); err != nil {
			// Fail open: a transient API error must not block valid rescales.
			// The next tick re-checks.
			log.Printf("rescaler: server %d: pre-check %q: %v (proceeding)", srv.ID, t, err)
		} else if !avail {
			if st != nil {
				now := time.Now().UTC()
				loc := srv.Datacenter.Location.Name
				_, _ = st.AppendEvent(store.Event{
					ServerID:    serverID,
					Kind:        "rescale_skipped",
					FromType:    t,
					ToType:      "",
					StartedAt:   now,
					FinishedAt:  now,
					OK:          false,
					Error:       fmt.Sprintf("unavailable in %s", loc),
					TriggeredBy: triggeredBy,
					Phase:       "pre_check",
				})
			}
			continue
		}
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