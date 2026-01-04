# =========================
# BUILD STAGE
# =========================
FROM golang:1.23-alpine AS builder

WORKDIR /app

# OS dependency
RUN apk add --no-cache git

# Cache dependency
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" \
    -o app ./cmd/api

# =========================
# RUNTIME STAGE
# =========================
FROM alpine:latest

# Security + cert + timezone
RUN apk add --no-cache ca-certificates tzdata \
 && adduser -D appuser

WORKDIR /app

# Copy binary only
COPY --from=builder /app/app .

# ENV default (staging safe)
ENV APP_PORT=8080
ENV TZ=Asia/Jakarta

EXPOSE 8080

USER appuser

# Healthcheck (optional tapi disarankan)
HEALTHCHECK --interval=30s --timeout=5s \
 CMD wget -qO- http://localhost:8080/health || exit 1

CMD ["./app"]
