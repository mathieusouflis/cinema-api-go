#!/bin/sh
set -e

MIGRATIONS_DIR="./db/migrations"

if [ -d "$MIGRATIONS_DIR" ] && [ -n "$(find "$MIGRATIONS_DIR" -name '*.sql' -print -quit 2>/dev/null)" ]; then
  echo "→ running migrations..."
  migrate -path "$MIGRATIONS_DIR" -database "$DATABASE_URL" up
  echo "→ migrations done"
fi

exec air -c .air.toml
