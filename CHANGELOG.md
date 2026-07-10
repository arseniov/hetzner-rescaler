# CHANGELOG

## Unreleased
- Run the per-server scheduler inside `serve` so scheduled / auto_promote
  events fire from the web UI without a separate `start` process. Multi-project
  support via `apiResolve`. Server create/update/delete now broadcast a
  `ServerLifecycleEvent` on the store's hub; both `serve` and `start`
  subscribe via `scheduler.Attach`.
- Pre-rescale availability gate: the rescaler now skips chain entries
  whose target server type is unavailable in the server's location,
  emitting a `rescale_skipped` event instead of shutting the server down
  (which previously caused a perpetual off → on cycle on sold-out types).
  The `GET /api/server-types` endpoint now requires `?location=X` and
  reports per-location availability; the web's server-type picker shows
  an "Unavailable" badge for sold-out types (which remain selectable —
  they may come back).
  - **Hetzner SDK v2 deprecation fix:** `liveServerState` now reads the
    canonical `Server.Location` field (the new "where is this server"
    field per Hetzner's 2025-12-16 phase-out of datacenters), preferring
    it over the deprecated `Server.Datacenter` field which is scheduled
    for removal after 1 July 2026. Without this fix, `/api/servers/[id]`
    silently dropped the `location` field for any server whose response
    populated only the new field — which broke the per-location catalog
    gate on the edit page. The earlier code only read the deprecated
    path; both fields are still read for the transitional window.
  - **Web frontend belt-and-braces:** the server edit and windows pages
    now unconditionally fire `serverTypes.load(FALLBACK_LOCATION)` on
    mount (before awaiting `/api/servers/[id]`), then refire with the
    real location if it arrives and differs from the fallback. This
    makes `/api/server-types?location=X` reliable even for servers
    whose live response omits `location` entirely. Tests cover the
    no-location, same-location, and different-location paths.
  - **Notes for operators:** this release upgrades `hcloud-go` from
    v1.33.0 to v2.44.0 (transitive type widening; any in-tree helpers
    that touched the SDK may need to adapt to v2's API surface) and
    raises the `go.mod` `go` directive to `1.25.0` — building from
    source on older Go toolchains will fail.

## 1.0.2
**2021-11-20**
- Added support for timezones
- Show info about time and timezone on startup
- Fixed No such file error on config creation
- Fixed copypasted sentence in top server prompt & typos

## 1.0.1
**2021-11-17**
- Added support for ARMv7 (release)
- Upgraded to cross-arch docker image, now supporting ARMv7, ARM64 and AMD64
- Fixed wrong machine type selector
- Fixed missing config file if provided by --config flag
- Added license, changelog and readme in release

## 1.0.0
**2021-11-16**
- First release