# =========================
# BUILD STAGE
# =========================
FROM golang:1.26.6-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates

# Download dependencies first untuk memanfaatkan Docker layer cache
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build static binary (CGO off agar binary portable & ringan)
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -o /app/bulky-api ./cmd/api

# =========================
# RUNTIME STAGE
# =========================
FROM alpine:3.20

WORKDIR /app

RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    ffmpeg \
    libwebp-tools \
    curl \
    --repository=https://dl-cdn.alpinelinux.org/alpine/edge/main \
    postgresql-client

ENV TZ=UTC
ENV APP_PORT=8080

COPY --from=builder /app/bulky-api /app/bulky-api

# Jalankan migration production sebelum API menerima traffic.
COPY docker-entrypoint.sh /app/docker-entrypoint.sh
RUN chmod 755 /app/docker-entrypoint.sh

# golang-migrate CLI agar migrasi bisa dijalankan dari dalam container
# (berguna saat DB hanya accessible di dalam Docker network)
COPY --from=migrate/migrate /migrate /usr/local/bin/migrate
# File migration di-copy ke image agar entrypoint dapat menjalankannya saat startup.
COPY migrations ./migrations

EXPOSE 8080

# Health check agar orchestrator (Dokploy) tahu kapan container siap menerima traffic
HEALTHCHECK --interval=30s --timeout=5s --start-period=15s --retries=3 \
  CMD curl -f http://localhost:8080/api/health || exit 1

ENTRYPOINT ["/app/docker-entrypoint.sh"]
CMD ["/app/bulky-api"]
