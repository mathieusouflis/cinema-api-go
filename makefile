# ── Config ────────────────────────────────────────────────────────────────────
SERVICES := gateway auth user catalog people social watchparty hub notification

# ── Dev ───────────────────────────────────────────────────────────────────────

.PHONY: infra/up
infra/up: ## Start local infrastructure (postgres, redis, nats, meilisearch, minio)
	docker compose -f deploy/docker-compose.infra.yml up -d

.PHONY: infra/down
infra/down: ## Stop local infrastructure
	docker compose -f deploy/docker-compose.infra.yml down

.PHONY: infra/logs
infra/logs: ## Show infrastructure logs
	docker compose -f deploy/docker-compose.infra.yml logs -f

# ── Build ─────────────────────────────────────────────────────────────────────

.PHONY: build
build: ## Build all services
	@for svc in $(SERVICES); do \
		echo "→ building $$svc..."; \
		(cd services/$$svc && go build -o ../../bin/$$svc ./cmd/); \
	done

.PHONY: build/%
build/%: ## Build a single service  (e.g. make build/auth)
	cd services/$* && go build -o ../../bin/$* ./cmd/

# ── Run ───────────────────────────────────────────────────────────────────────

.PHONY: run/%
run/%: ## Run a single service  (e.g. make run/auth)
	cd services/$* && go run ./cmd/

# ── Test ──────────────────────────────────────────────────────────────────────

.PHONY: test
test: ## Run all tests
	@for svc in $(SERVICES); do \
		(cd services/$$svc && go test ./...) 2>/dev/null || true; \
	done

.PHONY: test/%
test/%: ## Run tests for a single service  (e.g. make test/auth)
	cd services/$* && go test ./...

.PHONY: test/cover
test/cover: ## Run tests with coverage report
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# ── Lint ──────────────────────────────────────────────────────────────────────

.PHONY: lint
lint: ## Run golangci-lint on all services
	golangci-lint run ./...

.PHONY: fmt
fmt: ## Format all Go files
	gofmt -w .
	goimports -w .

# ── Database ──────────────────────────────────────────────────────────────────

# MIGRATE — wrapper around scripts/migrate/main.go (no system tool required).
# Reads DATABASE_URL from services/$(svc)/.env automatically.
MIGRATE := go run ./scripts/migrate/
ENV_FILE  = services/$(svc)/.env

.PHONY: migrate/up
migrate/up: ## Run all pending migrations  (e.g. make migrate/up svc=auth)
	@set -a && . $(ENV_FILE) && set +a && \
	MIGRATIONS_PATH=services/$(svc)/db/migrations $(MIGRATE) up

.PHONY: migrate/down
migrate/down: ## Rollback last migration  (e.g. make migrate/down svc=auth)
	@set -a && . $(ENV_FILE) && set +a && \
	MIGRATIONS_PATH=services/$(svc)/db/migrations $(MIGRATE) down

.PHONY: migrate/create
migrate/create: ## Create a new migration  (e.g. make migrate/create svc=auth name=add_sessions)
	@set -a && . $(ENV_FILE) && set +a && \
	MIGRATIONS_PATH=services/$(svc)/db/migrations $(MIGRATE) create $(name)

.PHONY: db/seed
db/seed: ## Seed the database with test data
	psql "$(DATABASE_URL)" -f db/seed/seed.sql

# ── Code generation ───────────────────────────────────────────────────────────

.PHONY: sqlc
sqlc: ## Generate Go code from SQL queries (all services)
	@for svc in $(SERVICES); do \
		if [ -f services/$$svc/sqlc.yaml ]; then \
			echo "→ sqlc $$svc..."; \
			cd services/$$svc && sqlc generate && cd ../..; \
		fi \
	done

.PHONY: mocks
mocks: ## Generate mocks from domain interfaces (all services)
	@for svc in $(SERVICES); do \
		if [ -d services/$$svc/internal/domain ]; then \
			echo "→ mocks $$svc..."; \
			mockery --dir services/$$svc/internal/domain \
			        --output services/$$svc/internal/mocks \
			        --all; \
		fi \
	done

# ── Docker ────────────────────────────────────────────────────────────────────

.PHONY: docker/build
docker/build: ## Build all Docker images
	@for svc in $(SERVICES); do \
		echo "→ docker build $$svc..."; \
		docker build -t kirona/$$svc:latest \
		             -f services/$$svc/Dockerfile .; \
	done

.PHONY: docker/build/%
docker/build/%: ## Build a single Docker image  (e.g. make docker/build/auth)
	docker build -t kirona/$*:latest -f services/$*/Dockerfile .

# ── Helpers ───────────────────────────────────────────────────────────────────

.PHONY: tidy
tidy: ## Run go mod tidy on all services
	@for svc in $(SERVICES); do \
		echo "→ tidy $$svc..."; \
		cd services/$$svc && go mod tidy && cd ../..; \
	done
	go work sync

.PHONY: help
help: ## Show this help
	@grep -E '^[a-zA-Z_/%]+:.*?## .*$$' $(MAKEFILE_LIST) \
	| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-25s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
