#!/bin/sh
# Usage: ./scripts/new-service.sh <name> <port> [has-db: yes|no]
# Example: ./scripts/new-service.sh catalog 8083 yes
set -e

NAME=$1
PORT=$2
HAS_DB=${3:-yes}

if [ -z "$NAME" ] || [ -z "$PORT" ]; then
  echo "Usage: $0 <name> <port> [has-db: yes|no]"
  exit 1
fi

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
SVC_DIR="$ROOT/services/$NAME"

if [ -d "$SVC_DIR" ]; then
  echo "✗ services/$NAME already exists"
  exit 1
fi

echo "→ scaffolding service: $NAME (port $PORT, has-db: $HAS_DB)"

# ── Directory structure ───────────────────────────────────────────────────────
mkdir -p "$SVC_DIR/cmd"
mkdir -p "$SVC_DIR/scripts"
if [ "$HAS_DB" = "yes" ]; then
  mkdir -p "$SVC_DIR/db/migrations"
fi

# ── go.mod ────────────────────────────────────────────────────────────────────
cat > "$SVC_DIR/go.mod" << EOF
module $NAME

go 1.26.1
EOF

# ── .env ─────────────────────────────────────────────────────────────────────
cat > "$SVC_DIR/.env" << EOF
PORT=$PORT
EOF

# ── cmd/main.go ──────────────────────────────────────────────────────────────
cat > "$SVC_DIR/cmd/main.go" << EOF
package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "$PORT"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "ok")
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello from $NAME service")
	})

	slog.Info("starting", "service", "$NAME", "port", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		slog.Error("server error", "err", err)
		os.Exit(1)
	}
}
EOF

# ── .air.toml ─────────────────────────────────────────────────────────────────
cat > "$SVC_DIR/.air.toml" << 'EOF'
root = "."
tmp_dir = "tmp"

[build]
  cmd = "go build -o ./tmp/main ./cmd/"
  bin = "./tmp/main"
  full_bin = "./tmp/main"
  include_ext = ["go"]
  exclude_dir = ["tmp", "vendor"]
  include_dir = [".", "../../pkg"]
  delay = 500
  stop_on_error = true
  kill_delay = "1s"

[log]
  time = false

[color]
  main = "magenta"
  watcher = "cyan"
  build = "yellow"
  runner = "green"

[misc]
  clean_on_exit = true
EOF

# ── Dockerfile.dev ────────────────────────────────────────────────────────────
if [ "$HAS_DB" = "yes" ]; then
  MIGRATE_INSTALL='    go install -tags '"'"'postgres'"'"' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.18.1'
  ENTRYPOINT='ENTRYPOINT ["sh", "scripts/entrypoint.dev.sh"]'
else
  MIGRATE_INSTALL=""
  ENTRYPOINT='CMD ["air", "-c", ".air.toml"]'
fi

cat > "$SVC_DIR/Dockerfile.dev" << EOF
FROM golang:1.26-alpine

RUN apk add --no-cache git && \\
    go install github.com/air-verse/air@latest${MIGRATE_INSTALL:+ && \\
$MIGRATE_INSTALL}

WORKDIR /workspace

COPY go.mod go.sum ./
COPY pkg/ ./pkg/
COPY services/$NAME/ ./services/$NAME/

RUN printf 'go 1.26.1\\\n\\\nuse (\\\n\\\t.\\\n\\\t./services/$NAME\\\n)\\\n' > go.work
RUN go work sync && cd services/$NAME && go mod download

WORKDIR /workspace/services/$NAME

$ENTRYPOINT
EOF

# ── entrypoint.dev.sh (only if has DB) ───────────────────────────────────────
if [ "$HAS_DB" = "yes" ]; then
  cat > "$SVC_DIR/scripts/entrypoint.dev.sh" << 'EOF'
#!/bin/sh
set -e

MIGRATIONS_DIR="./db/migrations"

if [ -d "$MIGRATIONS_DIR" ] && [ -n "$(find "$MIGRATIONS_DIR" -name '*.sql' -print -quit 2>/dev/null)" ]; then
  echo "→ running migrations..."
  migrate -path "$MIGRATIONS_DIR" -database "$DATABASE_URL" up
  echo "→ migrations done"
fi

exec air -c .air.toml
EOF
  chmod +x "$SVC_DIR/scripts/entrypoint.dev.sh"
fi

# ── go.work ───────────────────────────────────────────────────────────────────
# Insert the new service into go.work before the closing )
sed -i.bak "/)/{
  i\\
\\t./services/$NAME
}" "$ROOT/go.work" && rm -f "$ROOT/go.work.bak"

# ── docker-compose.dev.yml ────────────────────────────────────────────────────
# Append service block before the end of file
# NOTE: This is a best-effort append; review the compose file after running.
cat >> "$ROOT/deploy/docker-compose.dev.yml" << EOF

  # ── $(echo "$NAME" | tr '[:lower:]' '[:upper:]') ─────────────────────────────────────────────────────────
  $NAME:
    build:
      context: ..
      dockerfile: services/$NAME/Dockerfile.dev
    networks: [services]
    environment:
      PORT: "$PORT"
    env_file:
      - ../services/$NAME/.env
    develop:
      watch:
        - path: ../services/$NAME
          action: sync
          target: /workspace/services/$NAME
          ignore: [tmp, vendor]
        - path: ../services/$NAME/go.mod
          action: rebuild
        - path: ../pkg
          action: sync
          target: /workspace/pkg
EOF

echo "✓ scaffolded services/$NAME"
echo ""
echo "Next steps:"
echo "  1. Add networks (postgres, redis, nats...) to the service block in deploy/docker-compose.dev.yml"
echo "  2. Add environment variables (DATABASE_URL, etc.) to the service block"
if [ "$HAS_DB" = "yes" ]; then
  echo "  3. Add kirona_$NAME database to scripts/db/init.sql"
fi
echo "  4. Run: make infra/up && make dev/watch"
