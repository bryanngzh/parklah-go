#!/usr/bin/env bash

# Exit immediately if a command exits with a non-zero status.
set -e

# Load environment variables
source .env

# Check if DB_URL is set, else build it dynamically
if [ -z "$DB_URL" ]; then
  export DB_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}"
fi

echo "🚀 Running Goose migrations on: $DB_URL"

# Run Goose migrations
goose -dir ./db/migrations postgres "$DB_URL" up

echo "✅ Migration complete!"
