VERSION = ${shell git describe --tag}
ifeq (VERSION, "")
  VERSION = "v0-alpha"
endif

COMMAND := $(filter-out $(firstword $(MAKECMDGOALS)), $(MAKECMDGOALS))
FLAGS := $(filter-out -Werror, $(CFLAGS))

watch:
	gow run ./main.go $(COMMAND)

run:
	go run ./main.go $(COMMAND)

serve:
	./bin/hetzner-rescaler $(COMMAND)

# Build including dynamic libraries
build-static:
	CGO_ENABLED=0 && GOOS=linux && GOARCH=amd64 && go build -a -tags netgo -ldflags '-w -extldflags "-static"' -o ./bin/hetzner-rescaler_static ./main.go

# Multi-platform build requires buildx installed
build-multiarch-docker:
	go mod vendor && \
	docker buildx build --platform linux/arm/v7,linux/amd64,linux/arm64 -t jonamat/hetzner-rescaler:latest -t jonamat/hetzner-rescaler:${VERSION} . --push

build-docker:
	docker build -t jonamat/hetzner-rescaler:latest -t jonamat/hetzner-rescaler:${VERSION} .

push-docker:
	docker push jonamat/hetzner-rescaler:latest jonamat/hetzner-rescaler:${VERSION}

create-release:
	./scripts/release.sh
test-engine:
	go test -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Phase 2: build the SPA then the Go binary. The SPA is no longer
# embedded in the binary — production builds happen in Docker (web/
# and caddy/ each have their own Dockerfile).
#
# Note: this project does NOT vendor dependencies (vendor/ is gitignored),
# so we do not pass -mod vendor to go build. The module cache is used
# directly; run `go mod download` manually if modules are missing.

.PHONY: web-deps web-build serve-dev

web-deps:
	cd web && bun install

web-build: web-deps
	cd web && bun run build

build: web-build
	go build -o ./bin/hetzner-rescaler ./main.go

# Run the local Go binary against a local DB and loopback HTTP. The
# SPA is served separately via `cd web && PUBLIC_INTERNAL_TOKEN=dev-token
# bun run dev` (which uses Vite's /api proxy to talk to this process).
serve-dev:
	RESCALER_DB_PATH=$$HOME/.hetzner-rescaler/dev.db \
	RESCALER_INTERNAL_TOKEN=dev-token \
	RESCALER_HTTP_ADDR=127.0.0.1:8080 \
	./bin/hetzner-rescaler serve