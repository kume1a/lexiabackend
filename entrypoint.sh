#!/bin/sh
set -e

echo "Running database migrations..."
atlas migrate apply --dir "file://ent/migrate/migrations" --url "$DB_CONNECTION_URL"

echo "Starting application..."
exec /lexiabin
