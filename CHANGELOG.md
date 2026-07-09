# CHANGELOG

## Unreleased
- Run the per-server scheduler inside `serve` so scheduled / auto_promote
  events fire from the web UI without a separate `start` process. Multi-project
  support via `apiResolve`. Server create/update/delete now broadcast a
  `ServerLifecycleEvent` on the store's hub; both `serve` and `start`
  subscribe via `scheduler.Attach`.

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