FROM golang:1.23-bullseye AS builder
WORKDIR /build

ARG TARGETOS
ARG TARGETARCH

RUN apt update && apt install ca-certificates && apt install tzdata

COPY . .

# Create statically linked server binary
RUN CGO_ENABLED=0 && GOOS=${TARGETOS} && GOARCH=${TARGETARCH} && go build -x -mod vendor -tags netgo -ldflags '-w -extldflags "-static"' -o ./bin/hetzner-rescaler ./main.go

FROM scratch AS runner
WORKDIR /

COPY --from=builder /build/bin/hetzner-rescaler /bin/hetzner-rescaler
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Persistent state lives here in the container:
#   - /data/db.sqlite  (RESCALER_DB_PATH)
#   - /data/key         (encryption key, mode 0600)
VOLUME ["/data"]

# Sensible defaults — overridable via `docker run -e ...` or compose `environment:`.
# Phase 1: scheduler loop. Phase 2: switches to `serve` (loopback HTTP + embedded SPA).
ENV RESCALER_DB_PATH=/data/db.sqlite \
    RESCALER_HTTP_ADDR=127.0.0.1:8080

# By default we run the scheduler loop. Phase 2 will change this to `serve`.
# Operators can override with `docker run ... hetzner-rescaler <subcommand>`.
CMD ["hetzner-rescaler", "start"]
