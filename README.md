# Hetzner Rescaler

<p align="center">
  <a href="https://github.com/jonamat/hetzner-rescaler/releases">
    <img alt="Release" src="https://img.shields.io/github/v/release/jonamat/hetzner-rescaler" />
  </a>

  <a href="https://hub.docker.com/repository/docker/jonamat/hetzner-rescaler">
    <img alt="Docker Image Size (tag)" src="https://img.shields.io/docker/image-size/jonamat/hetzner-rescaler/latest" />
  </a>

  <a href="https://github.com/jonamat/hetzner-rescaler/blob/master/go.mod">
    <img alt="GitHub go.mod Go version" src="https://img.shields.io/github/go-mod/go-version/jonamat/hetzner-rescaler" />
  </a>

  <a href="https://github.com/jonamat/hetzner-rescaler/actions">
    <img alt="GitHub Workflow Status" src="https://github.com/jonamat/hetzner-rescaler/actions/workflows/push.yml/badge.svg" />
  </a>
</p>

Self-hosted scheduler that programmatically rescales your Hetzner Cloud
servers — dropping to a cheaper VM when load is low and promoting to a
more capable one when it's high — managed from a browser, with multiple
projects and per-project credentials.

## Features

- **One container, multiple Hetzner projects.** Each project gets its
  own API token (encrypted at rest), its own server inventory, and its
  own schedule.
- **Three rescale modes per server:** `scheduled` (one or more
  day/time windows), `auto_promote` (scale up under sustained load
  and back down), `manual` (operator-driven via UI or CLI).
- **Web UI** — SvelteKit + Better Auth SPA with a dashboard, projects
  page, servers list, server detail, windows editor, events timeline,
  status pages, and a live SSE stream of rescale activity.
- **CLI** — `hetzner-rescaler` for headless / CI use. Stores everything
  in SQLite; no config files to ship around.
- **Single shared database.** The Go engine and the SvelteKit SPA
  read/write the same SQLite file over WAL — no replication, no
  separate users, no Drizzle/admin tool needed.

## Quick start

You need at least one [Hetzner Cloud API
token](https://docs.hetzner.cloud/#getting-started) (read+write scope)
to import a project. Two install paths:

### Option A — Docker Compose (production / shared server)

```sh
cp .env.example .env                       # set RESCALER_INTERNAL_TOKEN + BETTER_AUTH_SECRET
docker compose up -d --build
open http://localhost:8089                 # sign up; the first account is the admin
```

Three services come up behind a single port (default 8089):

| Service        | Image                              | Role                                            |
| -------------- | ---------------------------------- | ----------------------------------------------- |
| `caddy`        | `jonamat/hetzner-rescaler-caddy`   | Public entrypoint; routes by path.             |
| `rescaler-api` | `jonamat/hetzner-rescaler`         | Go HTTP API + scheduler loop.                  |
| `rescaler-web` | `jonamat/hetzner-rescaler-web`     | SvelteKit SPA + Better Auth on Bun + bun:sqlite |

See [`compose.yaml`](./compose.yaml) for the full service config.

### Option B — Local development with Bun + Go (no Docker)

For hacking on the SPA, the scheduler, or the API in isolation.
You'll need [Go ≥ 1.23](https://go.dev/dl/) and
[Bun ≥ 1.1](https://bun.sh/).

```sh
# 1. Generate required secrets
export RESCALER_INTERNAL_TOKEN=$(head -c 32 /dev/urandom | xxd -p -c 64)
export BETTER_AUTH_SECRET=$(head -c 32 /dev/urandom | base64)

# 2. Build the Go CLI
make build                                 # outputs ./bin/hetzner-rescaler

# 3. Start the API + scheduler on :8080
make serve-dev

# 4. In another terminal — start the SPA on :5173 with HMR
cd web
PUBLIC_INTERNAL_TOKEN="$RESCALER_INTERNAL_TOKEN" bun run dev
open http://localhost:5173                 # Vite proxies /api/* to :8080
```

The Go binary serves `/api/*` (except `/api/auth/*`) on loopback. The
SPA's `vite.config.ts` proxies `/api/*` to `http://127.0.0.1:8080`,
which is the same loopback the Go process is bound to. Better Auth's
`/api/auth/*` calls are served from SvelteKit itself.

## Architecture

The runtime ships as three loosely-coupled services that share one
SQLite file via a Docker volume.

```
                  ┌─────────────────────────────────────────────┐
                  │            Browser (your laptop)            │
                  └────────────────────────┬────────────────────┘
                                           │   http://localhost:8089
                                           ▼
             ┌─────────────────────────────────────────────────────┐
             │                        caddy                        │
             │  ┌─────────────────┐         ┌─────────────────┐   │
             │  │  Host: …:8089   │         │  Host: …:8089   │   │
             │  │  /api/auth/*    │         │  /api/* (rest)  │   │
             │  │  /* (fallback)  │         │                 │   │
             │  └────────┬────────┘         └────────┬────────┘   │
             └───────────┼──────────────────────────┼─────────────┘
                         │                          │
                         ▼                          ▼
             ┌──────────────────────┐    ┌──────────────────────┐
             │     rescaler-web     │    │     rescaler-api     │
             │     (Bun runtime)    │    │      (Go binary)     │
             │                      │    │                      │
             │ SvelteKit SPA shell  │    │ net/http router      │
             │   + /api/auth/*      │    │   /api/healthz       │
             │     (Better Auth,    │    │   /api/projects      │
             │      Drizzle,        │    │   /api/servers       │
             │      bun:sqlite)     │    │   /api/windows       │
             │                      │    │   /api/events        │
             │ X-Internal-Token     │    │   /api/events/stream │
             │   middleware → API   │    │   /api/server-types  │
             └────────────┬─────────┘    │   /api/metrics       │
                          │              │                      │
                          │              │ scheduler goroutines │
                          │              │   tick → rescaler.   │
                          │              │   RescaleWithFallback│
                          │              └──────────┬───────────┘
                          │                         │
                          ▼                         ▼
                       ┌────────────────────────────────────────┐
                       │        SQLite (WAL mode, fs volume)    │
                       │                                        │
                       │   projects · servers · windows          │
                       │   actions · events · __drizzle_mig      │
                       └────────────────────────────────────────┘
```

**Why three services?**
- The Go engine is the single source of truth for rescale decisions
  and event records; the SPA is a stateless renderer.
- The SPA holds Better Auth's user/session/account tables in the same
  SQLite file, so user creation and rescale history live in the
  exact same backup.
- Caddy routes purely by path prefix. Both upstream services see
  the browser's `Host` header verbatim (`header_up Host {host}`),
  which keeps cookies, redirect URLs, and SvelteKit's
  `event.url` reconstruction working.

## Web UI

[Flowbite Svelte](https://flowbite-svelte.com/) + ApexCharts on top
of [SvelteKit 2](https://svelte.dev/) with [Better
Auth](https://www.better-auth.com/) for sign-up / sign-in. All data
is fetched from the Go API over loopback; the SPA never talks to
Hetzner directly.

| Page        | Purpose                                                            |
| ----------- | ------------------------------------------------------------------ |
| `/`         | Dashboard — KPIs, sparkline charts, and recent events.            |
| `/login`    | Better Auth sign-up / sign-in. First account becomes admin.       |
| `/projects` | List / add / remove Hetzner projects; per-project refresh button. |
| `/servers`  | All servers across projects; mode + windows badges.                |
| `/servers/:id` | Per-server detail: edit mode, manage rescale windows, view events. |
| `/events`   | Global event log (rescales, refreshes, errors) with filters.      |
| `/status/servers`, `/status/health` | Liveness pages for the engine and per-server reachability. |

A live SSE connection at `/api/events/stream` pushes every event into
the dashboard the moment the scheduler commits it.

## CLI

`hetzner-rescaler` is the Go binary. It has two ways to run:

| Command   | What it does                                                              |
| --------- | ------------------------------------------------------------------------- |
| `serve`   | HTTP API + scheduler loop. **Required** for the SPA.                       |
| `start`   | Scheduler loop only. No HTTP, no UI. Useful for headless / CI / VM-only deploys. |
| `config`  | Interactive REPL — add/edit projects, servers, modes, windows.            |
| `status`  | Print configured projects + servers + the most recent events to stdout.   |
| `try`     | One-shot rescale: `hetzner-rescaler try <server-id> <up\|down>`.          |
| `migrate` | Import a legacy v1 YAML config into the SQLite database.                  |

Generated config and rescale history live in a single SQLite file;
Hetzner API tokens are sealed with AES-256-GCM using
`RESCALER_TOKEN_ENCRYPTION_KEY`. The key is **required** — there is no
auto-generated fallback. A fresh container without the env var refuses
to start; this is intentional, because silently minting a new key on
image rebuild would lock out previously-stored tokens (GCM auth tag
mismatch). Generate one with `openssl rand -hex 32` and persist it
alongside your database backups.

### Multi-project architecture

Each row in `projects` owns its own Hetzner API token (encrypted at
rest) and its own inventory of servers. The scheduler instantiates one
goroutine per server and asks `store.Scheduler.apiFor(projectID)` for
the right `hetzner.API` client at every tick. Refresh on a project
imports its current server list in one HTTP call.

Three rescale modes per server:
- **`manual`** — operator-only. Buttons on the server detail page or
  `hetzner-rescaler try <id> up|down`.
- **`scheduled`** — one or more `windows` (day mask + start/stop time
  + target server type). On every tick, the scheduler computes the
  current window target and rescales if it differs from the running
  type.
- **`auto_promote`** — scale up under sustained CPU pressure and back
  down when the load eases.

## Configuration

All configuration is via environment variables. The bundled
[`.env.example`](./.env.example) lists every key with comments.

| Variable                       | Required            | Purpose                                                  |
| ------------------------------ | ------------------- | -------------------------------------------------------- |
| `RESCALER_INTERNAL_TOKEN`      | for `serve`         | Shared secret between SPA and Go API (`X-Internal-Token`). Baked into the SPA at build time. |
| `RESCALER_TOKEN_ENCRYPTION_KEY`| **required**        | Hex-encoded 32-byte AES-GCM key (64 hex chars). Generate with `openssl rand -hex 32`. Without it the API refuses to start. |
| `RESCALER_HTTP_ADDR`           | for `serve`         | Listen address for the API (default `0.0.0.0:8080`).     |
| `RESCALER_DB_PATH`             | recommended         | SQLite path (default `./db.sqlite`, `/data/db.sqlite` in Docker). |
| `BETTER_AUTH_SECRET`           | for `serve`         | Signs Better Auth session tokens; ≥ 32 chars.            |
| `BETTER_AUTH_URL`              | for `serve`         | Public origin the browser uses; used for cookie scoping. |
| `DATABASE_URL`                 | for `serve`         | SQLite path for the SPA. Must equal `RESCALER_DB_PATH`.  |
| `PUBLIC_INTERNAL_TOKEN`        | SPA build arg       | Mirror of `RESCALER_INTERNAL_TOKEN`; baked into the client bundle. |
| `ORIGIN`                       | SPA runtime         | Tells adapter-node the request scheme+origin. Must match `BETTER_AUTH_URL`. |
| `TZ`                           | optional            | Timezone for window evaluation.                          |

## Backup

The SQLite database **and** the encryption key must be backed up
together — backup one without the other and your Hetzner tokens are
unrecoverable.

| Component | Local dev                | Docker                          |
| --------- | ------------------------ | ------------------------------- |
| DB        | `~/.hetzner-rescaler/db.sqlite` (or `$RESCALER_DB_PATH`) | `/data/db.sqlite` |
| Key       | `$RESCALER_TOKEN_ENCRYPTION_KEY` env (no on-disk fallback anymore) | `$RESCALER_TOKEN_ENCRYPTION_KEY` env |

Back up the entire `/data/` volume.

## Migration from v1 (YAML)

If you have a v1 `~/.hetzner-rescaler.yaml` lying around, convert it
in-place:

```sh
hetzner-rescaler migrate            # uses the default path
hetzner-rescaler migrate --from /path/to/old.yaml
```

The command creates a single project with the configured server and
preserves the original schedule.

## Development

### Building from source

```sh
# Go engine
make build                             # → ./bin/hetzner-rescaler

# SPA (bun + bun:sqlite)
cd web && bun install && bun run build # → web/build/

# Multi-arch Docker images (requires buildx)
make build-multiarch-docker            # builds + pushes linux/amd64, arm/v7, arm64
```

### Tests

```sh
# Go (race detector + coverage report)
make test-engine

# SPA (vitest)
cd web && bun run test
```

### Project layout

```
.
├── cmd/                  cobra subcommands: config, serve, start, status, try, migrate
├── internal/
│   ├── api/              HTTP routes, SSE, auth middleware, handler tests
│   ├── broadcast/        generic in-process pub/sub Hub[T] for live events
│   ├── crypto/           AES-256-GCM keyring for Hetzner tokens
│   ├── hcloudmock/       fake Hetzner SDK for handler tests
│   ├── hetzner/          typed wrapper around the official Hetzner SDK
│   ├── rescaler/         RescaleWithFallback (the one path used by both API + scheduler)
│   ├── scheduler/        per-server tick goroutines + window/time evaluation
│   └── store/            the only package that knows SQL (migrations, CRUD, events)
├── web/                  SvelteKit 2 + Svelte 5 + Bun + bun:sqlite + Better Auth
├── caddy/                Caddy Dockerfile + Caddyfile (path-based reverse proxy)
├── compose.yaml          single-port service graph
└── main.go               thin entrypoint → cmd.Execute()
```

### Adding a new rescale trigger

The single dispatch path is `internal/rescaler.RescaleOnce`. Both
the API handlers (`POST /api/servers/:id/rescale|promote|demote`)
and the scheduler tick call it. New triggers (cron, webhook,
manual CLI) should wrap this same function so events end up in the
same `events` table and the SPA's SSE stream picks them up.

## License

MIT — see [LICENSE](./LICENSE).
