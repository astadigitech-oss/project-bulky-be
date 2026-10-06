#!/bin/sh
set -eu

if [ "${APP_ENV:-development}" = "production" ]; then
  : "${DB_USER:?DB_USER must be set for production migrations}"
  : "${DB_PASSWORD:?DB_PASSWORD must be set for production migrations}"
  : "${DB_HOST:?DB_HOST must be set for production migrations}"
  : "${DB_PORT:?DB_PORT must be set for production migrations}"
  : "${DB_NAME:?DB_NAME must be set for production migrations}"

  database_url="postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable"

  echo "Applying database migrations..."
  migrate -path /app/migrations -database "$database_url" up
  echo "Database migrations complete."
else
  echo "Skipping automatic database migrations because APP_ENV is not production."
fi

exec "$@"
