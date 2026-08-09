# ============================
# BUILD STAGE
# ============================
FROM golang:1.26.5-alpine3.23 AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates tzdata

# Copy go.mod/go.sum dulu supaya layer cache efektif saat dependency tidak berubah
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Info versi di-embed ke binary lewat build-arg (isi dari CI/CD pipeline)
ARG VERSION=dev
ARG COMMIT=unknown
ARG BUILD_DATE=unknown

# CGO_ENABLED=0        -> binary statis, tidak butuh libc dinamis
# -trimpath            -> hapus path lokal dari binary (reproducible build)
# -ldflags="-s -w"     -> strip simbol debug, ukuran binary lebih kecil
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -trimpath \
    -ldflags="-s -w -X main.Version=${VERSION} -X main.Commit=${COMMIT} -X main.BuildDate=${BUILD_DATE}" \
    -o /app/bin/api ./cmd/api

# ============================
# RUNTIME STAGE
# ============================
FROM alpine:3.23.2

# Pin versi package sebisa mungkin di CI (apk add pkg=version) untuk build yang reproducible.
# tini dipakai sebagai PID 1 supaya SIGTERM/SIGINT diteruskan dengan benar (graceful shutdown).
RUN apk add --no-cache \
        bash \
        curl \
        postgresql-client \
        ffmpeg \
        ca-certificates \
        tzdata \
        tini \
    && addgroup -S appgroup -g 10001 \
    && adduser -S appuser -G appgroup -u 10001 -h /app -s /sbin/nologin

WORKDIR /app

# --chown langsung saat copy, tidak perlu RUN chown terpisah (hemat layer)
COPY --from=builder --chown=appuser:appgroup /app/bin/api /app/api

ENV APP_PORT=8080
ENV TZ=UTC

# Jalankan sebagai non-root
USER appuser

EXPOSE 8080

# Health check memakai endpoint /health yang sudah ada di aplikasi
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD curl -fsS "http://localhost:${APP_PORT}/health" || exit 1

# tini sebagai init process (PID 1) untuk signal handling & reap zombie process
ENTRYPOINT ["/sbin/tini", "--"]
CMD ["/app/api"]
