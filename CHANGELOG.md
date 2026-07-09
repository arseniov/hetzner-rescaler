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