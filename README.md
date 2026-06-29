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
  serve       Run the HTTP API server (phase 2 — currently a stub)
  start       Start the scheduler loop (reads from SQLite)
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

## Docker

A pre-built image is published as `jonamat/hetzner-rescaler:latest` (multi-arch: amd64, arm64, armv7).

```sh
cp .env.example .env                       # edit values
cp docker-compose.example.yml docker-compose.yml
docker compose up -d
docker compose exec rescaler hetzner-rescaler config   # one-time interactive setup
docker compose logs -f rescaler
```

The container stores its SQLite DB and encryption key under `/data`, backed by a named volume (`rescaler_data`). Override behavior via `.env`:

- `RESCALER_DB_PATH` — defaults to `/data/db.sqlite`
- `RESCALER_TOKEN_ENCRYPTION_KEY` — hex-encoded 32-byte key; auto-generated on first run if unset (and written to `/data/key`)
- `RESCALER_HTTP_ADDR` — loopback HTTP address (only relevant once phase 2's `serve` command is in use)

The web UI (phase 2) and the Authorizer sibling container are not yet wired in this image.

## Use cases
This tool was developed for a my specific use case: I use an Hetzner server for remote development, using the [Remote SSH extension](https://code.visualstudio.com/docs/remote/ssh) to simplify my cross-device development workflow. This machine also serve some personal services, which require very little resources but cannot be stopped for a long time.<br>
It could be useful for all servers running applications related to a company's opening hours, such as booking, delivery or management software.

## License
MIT
