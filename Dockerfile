# =========================
# BUILD STAGE
# =========================
FROM golang:1.26-alpine AS builder

WORKDIR /app

RUN apk add --no-cache \
    git \
    bash \
    curl \
    postgresql-client \
    tzdata

COPY . .

RUN go mod download

ENV APP_PORT=8080
ENV TZ=Asia/Jakarta

EXPOSE 8080

CMD ["go", "run", "cmd/api/main.go"]
