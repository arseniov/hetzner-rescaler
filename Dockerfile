# ---- Stage 1: build the SvelteKit SPA ----
FROM oven/bun:1.1 AS spa
WORKDIR /spa
COPY web/package.json web/bun.lock* ./
RUN bun install --frozen-lockfile || bun install
COPY web/ ./
ARG PUBLIC_INTERNAL_TOKEN=docker-placeholder
ARG PUBLIC_AUTHORIZER_URL=http://localhost:8080
ENV PUBLIC_INTERNAL_TOKEN=$PUBLIC_INTERNAL_TOKEN
ENV PUBLIC_AUTHORIZER_URL=$PUBLIC_AUTHORIZER_URL
RUN bun run build

# ---- Stage 2: build the Go binary with the SPA embedded ----
FROM golang:1.23-bullseye AS builder
WORKDIR /build

ARG TARGETOS
ARG TARGETARCH

RUN apt update && apt install -y ca-certificates tzdata

COPY . .
COPY --from=spa /internal/web/build ./internal/web/build

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -mod vendor -tags netgo \
    -ldflags '-w -extldflags "-static"' \
    -o ./bin/hetzner-rescaler ./main.go

# ---- Stage 3: scratch runner ----
FROM scratch AS runner
WORKDIR /

COPY --from=builder /build/bin/hetzner-rescaler /bin/hetzner-rescaler
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

VOLUME ["/data"]

ENV RESCALER_DB_PATH=/data/db.sqlite \
    RESCALER_HTTP_ADDR=127.0.0.1:8080

CMD ["hetzner-rescaler", "serve"]