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

build:
	go build -v -x -o ./bin/hetzner-rescaler ./main.go

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

build-spa:
	cd web && bun install && bun run build
	cp -R web/build ./web_build_artifact

# Phase 2: build the SPA first, then embed it in the Go binary.
# `bun run build` requires `bun install` to have been run at least once;
# the web-deps / web-build targets below handle that.
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
	CGO_ENABLED=0 GOOS=$$(go env GOOS) GOARCH=$$(go env GOARCH) \
	  go build -tags netgo \
	  -ldflags '-w -extldflags "-static"' \
	  -o ./bin/hetzner-rescaler ./main.go

# Run the embedded Go binary against a local DB and loopback HTTP.
# Requires `web/build/` to exist (run `make build` first).
serve-dev:
	RESCALER_DB_PATH=$$HOME/.hetzner-rescaler/dev.db \
	RESCALER_INTERNAL_TOKEN=dev-token \
	RESCALER_HTTP_ADDR=127.0.0.1:8080 \
	./bin/hetzner-rescaler serve
