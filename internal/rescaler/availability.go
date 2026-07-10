package rescaler

import (
	"context"

	"github.com/jonamat/hetzner-rescaler/internal/hetzner"
)

// IsTypeAvailable reports whether the given server type is available
// in the server's current datacenter location. Returns (false, nil)
// when:
//   - srv is nil, or
//   - srv.Datacenter is nil (e.g. a freshly-provisioned server that
//     hasn't been observed yet), or
//   - srv.Datacenter.Location is nil, or
//   - the type's Locations list has no entry for the requested
//     location, or
//   - the type has an entry for the location but Available=false.
//
// Returns (false, err) on a transport / API error so the caller can
// decide whether to fail open or fail closed (the rescaler fails open).
//
// This helper performs a single GetServerType call per invocation.
// The rescaler typically calls it once per chain entry (1–3 entries
// per rescale) which is well within Hetzner's free-tier quota.
func IsTypeAvailable(ctx context.Context, api hetzner.API, srv *hetzner.Server, typeName string) (bool, error) {
	if srv == nil || srv.Datacenter == nil || srv.Datacenter.Location == nil {
		return false, nil
	}
	loc := srv.Datacenter.Location.Name
	t, err := api.GetServerType(ctx, typeName)
	if err != nil {
		return false, err
	}
	for _, l := range t.Locations {
		if l.Location != nil && l.Location.Name == loc {
			return l.Available, nil
		}
	}
	return false, nil
}
