# ---- Stage 1: build the Go API ----
FROM golang:1.23-bullseye AS builder
WORKDIR /build

ARG TARGETOS
ARG TARGETARCH

RUN apt update && apt install -y ca-certificates tzdata

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -tags netgo \
    -ldflags '-w -extldflags "-static"' \
    -o ./bin/hetzner-rescaler ./main.go

# ---- Stage 2: scratch runner ----
FROM scratch AS runner
WORKDIR /

COPY --from=builder /build/bin/hetzner-rescaler /bin/hetzner-rescaler
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

VOLUME ["/data"]

ENV RESCALER_DB_PATH=/data/db.sqlite \
    RESCALER_HTTP_ADDR=0.0.0.0:8080 \
    RESCALER_INTERNAL_TOKEN="" \
    BETTER_AUTH_URL=""

EXPOSE 8080

CMD ["hetzner-rescaler", "serve"]
