# Hetzner rescaler

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

Lightweight CLI tool to programmatically rescale your Hetzner virtual server daily to optimize your budget spending, scaling to a cheaper machine when you don't need or need few resources, and scaling to a more performant one when you know the load will be higher.

## Usage 
First, you need to generate an [Hetzner API Token](https://docs.hetzner.cloud/#getting-started).<br> 

Next you need to create your configuration file or export the required environment variables for the tool.<br>
The `config` command helps you generate a valid configuration, warning you if there are any logic errors and validating your input.
```sh
hetzner-rescaler config
``` 

Keep in mind if env vars are defined, they will take priority.<br>
After setting the configuration, you can start the tool by running
```sh
hetzner-rescaler start
```

### Use with environmental variables
Export these env vars to override or completely bypass the generated configuration:
| Variable           | Description                                                                  |
| ------------------ | ---------------------------------------------------------------------------- |
| `HCLOUD_TOKEN`     | A valid [Hetzner API Token](https://docs.hetzner.cloud/#getting-started)<br> |
| `SERVER_ID`        | The ID of the target server<br>                                              |
| `BASE_SERVER_NAME` | The code of the cheap server type<br>                                        |
| `TOP_SERVER_NAME`  | The code of the high performance server type<br>                             |
| `HOUR_START`       | 24h format, colon separated hour when the server should be upgraded<br>      |
| `HOUR_STOP`        | 24h format, colon separated hour when the server should be downgraded<br>    |
| `TZ`               | If defined, change the timezone of the timer<br>                             |

### Use with Docker
Pull the image from dockerhub
```sh
docker pull jonamat/hetzner-rescaler
```

**Opt A:** Create a config file inside the container & start immediately *beta* 
```sh
docker run -ti jonamat/hetzner-rescaler hetzner-rescaler plug
```

**Opt B:** Mounting a configuration file 
```sh
docker run -v ~/.hetzner-rescaler.yaml:/.hetzner-rescaler.yaml jonamat/hetzner-rescaler
```

**Opt C:** Passing config as env vars 
```sh
docker run \
-e HCLOUD_TOKEN=abc123 \
-e SERVER_ID=4567 \
-e BASE_SERVER_NAME=cpx11 \
-e TOP_SERVER_NAME=cpx21 \
-e HOUR_START=09:00 \
-e HOUR_STOP=20:00 \
jonamat/hetzner-rescaler
```

You can also pass a partial configuration file and define the missing vars as env vars (useful eg to hide the API key) 

### Use with compose/swarm stacks
```yml
version: '3.7'

services:
  hetzner-rescaler:
    image: jonamat/hetzner-rescaler

    // Provide the env vars
    environment:
      HCLOUD_TOKEN: abc123
      SERVER_ID: 4567
      BASE_SERVER_NAME: cpx11
      TOP_SERVER_NAME: cpx21
      HOUR_START: "09:00"
      HOUR_STOP: "20:00"
    
    // ...or mount the config file
    volumes:
      - /var/hetzner/config-file.yaml:/.hetzner-rescaler.yaml
```

## The configuration file
The default path for the config file is `~/.hetzner-rescaler.yaml`.<br>
You can provide (and create) a custom config path passing the `--config /custom/path/config.yml` flag.<br>

Config yaml file example
```yaml
hcloud_token: abc123
server_id: 15393230
base_server_name: cx11
top_server_name: cx21
hour_start: "09:00"
hour_stop: "20:00"
```

## Commands
```
Usage:
  hetzner-rescaler [command]

Available Commands:
  config      Interactively add or edit projects, servers, modes, and windows (stored in SQLite)
  help        Help about any command
  migrate     Import a legacy YAML config into the SQLite database
  serve       Run the HTTP API + static SPA + scheduler (loopback)
  start       Run the scheduler loop only (no HTTP)
  status      Print all configured projects, servers, and recent events
  try         One-shot rescale: `hetzner-rescaler try <server-id> <up|down>`

Flags:
  -h, --help   help for hetzner-rescaler

Use "hetzner-rescaler [command] --help" for more information about a command.
```

## Backup

The SQLite database and the token-encryption key must be backed up together:

- DB: `~/.hetzner-rescaler/db.sqlite` (or wherever `RESCALER_DB_PATH` points, default `/data/db.sqlite` in Docker)
- Key: `~/.hetzner-rescaler/key` (or `/data/key` in Docker)

If you back up the DB without the key, your Hetzner tokens are unrecoverable.

## Web UI (phase 2)

The SvelteKit web UI is shipped as a static SPA embedded in the same Docker image. The browser talks to the Go backend over loopback HTTP using an `X-Internal-Token` shared secret; authentication for the SPA itself is delegated to a sibling [Authorizer](https://github.com/authorizerdev/authorizer) container.

### Quick start with docker compose

```sh
cp .env.example .env                       # edit values (especially RESCALER_INTERNAL_TOKEN)
cp docker-compose.example.yml docker-compose.yml
docker compose up -d --build
```

This brings up two services:
- `rescaler` on http://localhost:8081 (HTTP API + SPA)
- `authorizer` on http://localhost:8080 (admin UI; the SPA uses it for login)

### First-time setup

1. Visit http://localhost:8080 to open the Authorizer admin UI and create a user.
2. Visit http://localhost:8081. You will be redirected to `/login`.
3. Sign in with the credentials you just created.
4. In the dashboard, click **Projects → Add project**. Enter a name and a Hetzner Cloud API token.
5. Click **Refresh from Hetzner** to import your existing servers.
6. Click a server, then **Rescale up / down** to test the action. Check **Events** for results.

### Configuration

The full set of environment variables is in `.env.example`. The two required for the web UI are:

| Variable | Purpose |
|----------|---------|
| `RESCALER_INTERNAL_TOKEN` | Shared secret between the SPA and the Go backend. Generated once, baked into the SPA at build time. |
| `AUTHORIZER_URL` | URL the browser uses to reach the Authorizer container. |
| `AUTHORIZER_CLIENT_ID` | Must match the Authorizer container's `--client-id`. |

### Running the SPA in development

```sh
# Terminal 1 — Go backend
make serve-dev

# Terminal 2 — SPA with HMR against the Go backend
cd web && bun install && bun run dev
```

The Vite dev server proxies `/api/*` to `http://127.0.0.1:8080`, so login + project management work end-to-end without rebuilding the SPA.

### CORS

The Authorizer container's `--allowed-origins` must include the URL the browser uses to reach the rescaler (default `http://localhost:8081`). The bundled `docker-compose.example.yml` sets this. If you put the rescaler behind a reverse proxy on a different host, update `--allowed-origins` accordingly.

## Docker

The image runs `serve` by default (loopback HTTP + SPA + scheduler). To run the scheduler-only CLI loop (no HTTP, no UI) start the container with an explicit subcommand:

```sh
docker run -d --name rescaler -v rescaler_data:/data \
  -e RESCALER_INTERNAL_TOKEN=changeme \
  jonamat/hetzner-rescaler start
```

## Use cases
This tool was developed for a my specific use case: I use an Hetzner server for remote development, using the [Remote SSH extension](https://code.visualstudio.com/docs/remote/ssh) to simplify my cross-device development workflow. This machine also serve some personal services, which require very little resources but cannot be stopped for a long time.<br>
It could be useful for all servers running applications related to a company's opening hours, such as booking, delivery or management software.

## License
MIT
