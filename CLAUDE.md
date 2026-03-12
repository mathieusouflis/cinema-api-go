# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**Kirona** — a community cinema & series application with a Go microservices backend. The project uses a Go workspace (go.work) monorepo structure with 9 planned microservices. Only `auth` is actively implemented and registered in `go.work`.

**Planned build order:** Auth → Catalog (films only first) → Frontend skeleton → Social → User/Friends → Watchparty → Realtime Hub → Chrome Extension.

## Commands

```bash
# Infrastructure
make infra/up          # Start Docker infra (postgres, redis, nats, meilisearch, minio)
make infra/down        # Stop infrastructure
make infra/logs        # Follow infrastructure logs

# Development
make build             # Build all services to bin/
make build/auth        # Build a single service
make run/auth          # Run a single service
make test              # Run all tests with coverage
make test/auth         # Test a single service
make test/cover        # Generate coverage report

# Code quality
make lint              # golangci-lint
make fmt               # gofmt + goimports
make tidy              # go mod tidy all services

# Database (requires svc= parameter)
make migrate/up svc=auth                      # Run all pending migrations
make migrate/down svc=auth                    # Roll back last migration
make migrate/create svc=auth name=<name>      # Create new migration file
make sqlc                                     # Regenerate sqlc query code

# Code generation
make mocks             # Generate mockery mocks from domain interfaces
make docker/build      # Build all Docker images
make docker/build/auth # Build single service image
```

## Architecture

### Monorepo Structure

- `pkg/` — shared packages imported by all services (`filmserver` module)
- `services/<name>/` — each service is its own Go module with its own `go.mod`
- `deploy/` — Docker Compose files
- `scripts/db/init.sql` — initializes per-service PostgreSQL databases (e.g., `kirona_auth`)

New services must be added to `go.work` with a `use ./services/<name>` entry.

### 4-Layer Service Architecture

Every service follows this strict layering (enforced by import rules):

```
handler/      ← HTTP only: parse request, call usecase, write response
usecase/      ← Business logic: orchestration, error mapping
repository/   ← Data access: SQL (sqlc), Redis, external APIs
domain/       ← Models and interfaces — zero internal imports
```

**Import rules:**
- `handler` imports `usecase` only, never `repository` directly
- `usecase` depends on `domain` interfaces, never concrete repository implementations
- `repository` implements `domain` interfaces
- `domain` has no internal dependencies

### Shared Packages (`pkg/`)

| Package | Purpose |
|---|---|
| `pkg/config` | `config.Load(cfg any)` — viper-based env var loading into typed structs |
| `pkg/config` | `config.Required(fields map[string]string)` — panics on missing required vars |
| `pkg/config` | Embeddable config types: `Base`, `Postgres`, `Redis`, `NATS`, `S3`, `Meilisearch`, `JWT`, `TMDB` |
| `pkg/errors` | Sentinel domain errors (`ErrNotFound`, `ErrConflict`, `ErrForbidden`, `ErrBadRequest`, `ErrUnauth`) |
| `pkg/errors` | `errors.Render(w, err)` — maps domain errors to HTTP status codes |
| `pkg/jwt` | `SignToken()`, `VerifyToken()`, `ParseToken()`, `JWTClaims` struct — HS256 |
| `pkg/logger` | `logger.New(env)` — returns `*slog.Logger` (JSON in prod, text in dev) |
| `pkg/render` | `render.JSON`, `render.Created`, `render.NoContent` — HTTP response helpers |
| `pkg/server` | `server.Run(addr, handler, log)` — starts HTTP server, blocks until SIGINT/SIGTERM, graceful shutdown (10 s) |

Stub packages not yet implemented: `pkg/middleware`, `pkg/health`, `pkg/otel`, `pkg/paginate`.

### Infrastructure (Docker Compose)

- PostgreSQL 16 (port 5432) — one database per service, initialized by `scripts/db/init.sql`
- Redis 7 (port 6379) — cache, sessions, pub-sub
- NATS 2 with JetStream (port 4222) — async events between services
- Meilisearch v1.6 (port 7700) — full-text search
- MinIO (ports 9000/9001) — S3-compatible file storage, default bucket: `kirona-media`

### Communication Patterns

- Client → Gateway: HTTPS REST
- Client → Hub service: WebSocket
- Service → Service: HTTP (internal)
- Services → NATS: async pub/sub events
- Services → Redis: cache, sessions, pub-sub
- Services → PostgreSQL: pgx pool (sqlc-generated queries)

### Config Pattern

Each service defines its own config struct embedding shared types:

```go
type Config struct {
    config.Base
    config.Postgres
    config.JWT
    // add only what the service needs
}

func Load() *Config {
    cfg := &Config{}
    if err := config.Load(cfg); err != nil {
        panic(err)
    }
    config.Required(map[string]string{"JWT_SECRET": cfg.JWT.Secret})
    return cfg
}
```

### Handler Pattern

Each handler is a struct implementing `http.Handler`:

```go
type LoginHandler struct { uc *usecase.Usecase }

func (h *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    var req usecase.Input
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        errors.Render(w, errors.ErrBadRequest)
        return
    }
    out, err := h.uc.Execute(r.Context(), req)
    if err != nil { errors.Render(w, err); return }
    render.JSON(w, http.StatusOK, out)
}
```

Routes are registered in `handler/routes.go` via a `Mount(router, deps)` function. Dependency injection is wired in `handler/dependencies.go`.

### Auth Service — Current State

Implemented: `POST /auth/register`, `/auth/login`, `/auth/refresh`, `/auth/logout`.

- Refresh and logout use an HttpOnly secure `refresh_token` cookie.
- `services/auth/pkg/hash/` — Argon2id password hashing (service-local package, not in `pkg/`).
- `db/orm/` — sqlc-generated code from `db/schema.sql` + `db/queries/`.
- `db/migrations/` — directory exists but is empty; migrations not yet created.
- OAuth2 handler (`handler/oauth.go`) and usecase (`usecase/oauth/callback.go`) are stubs — not implemented. OAuth env vars (`OAUTH_GOOGLE_*`, `OAUTH_GITHUB_*`) are defined in `.env`.

## Tech Stack

- **Language**: Go 1.26.1
- **HTTP Router**: `go-chi/chi/v5`
- **Config**: `spf13/viper` (env vars only, no config files)
- **Logging**: `log/slog` (stdlib)
- **DB Driver**: `pgx`, queries generated by `sqlc`
- **Migrations**: `golang-migrate`
- **Password hashing**: Argon2id (`golang.org/x/crypto`)
- **Mocks**: `mockery`
- **Linting**: `golangci-lint`
