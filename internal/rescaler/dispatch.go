// Package rescaler exposes a one-shot RescaleWithFallback wrapper that
// also writes a single event row. Used by cmd/serve.go so the API and
// the scheduler share the same dispatch path.
package rescaler

import (
	"context"
	"time"

	"github.com/jonamat/hetzner-rescaler/internal/hetzner"
	"github.com/jonamat/hetzner-rescaler/internal/store"
)

// RescaleOnce runs RescaleWithFallback for the given (project, server,
// target) tuple and writes one event row capturing the outcome. The
// caller passes the API for the server's project.
//
// Errors are returned to the caller for logging; the event row is still
// appended even on failure.
func RescaleOnce(ctx context.Context, api hetzner.API, srv *store.Server, target string, s *store.Store) error {
	hserver, err := api.GetServer(ctx, srv.HCloudServerID)
	if err != nil {
		return err
	}
	fromType := ""
	if hserver != nil && hserver.ServerType != nil {
		fromType = hserver.ServerType.Name
	}
	now := time.Now().UTC()
	used, err := RescaleWithFallback(ctx, api, hserver, target, srv.FallbackChain)
	if err != nil {
		_, _ = s.AppendEvent(store.Event{
			ServerID:    srv.ID,
			Kind:        "rescale_failed",
			FromType:    fromType,
			ToType:      target,
			StartedAt:   now,
			FinishedAt:  now,
			OK:          false,
			Error:       err.Error(),
			TriggeredBy: "api",
		})
		return err
	}
	_, _ = s.AppendEvent(store.Event{
		ServerID:    srv.ID,
		Kind:        eventKindFor(srv, used),
		FromType:    fromType,
		ToType:      used,
		StartedAt:   now,
		FinishedAt:  now,
		OK:          true,
		TriggeredBy: "api",
	})
	return nil
}

// eventKindFor returns "rescale_up" when scaling up to top, "rescale_down"
// when scaling down to base, "rescale_fallback" otherwise.
func eventKindFor(srv *store.Server, toType string) string {
	if toType == srv.BaseServerType {
		return "rescale_down"
	}
	if toType == srv.TopServerType {
		return "rescale_up"
	}
	return "rescale_fallback"
}