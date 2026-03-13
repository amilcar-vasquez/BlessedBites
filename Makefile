# ====================================
# Project: BlessedBites
# Monolith SSR app  +  Go REST API  +  SvelteKit frontend
# ====================================

# Load environment variables from .env if it exists
ifneq (,$(wildcard .env))
	include .env
	export
endif

# ------------------------------------
# Tuneable defaults (override via .env or CLI)
# ------------------------------------
APP_ENV     ?= development
ADDR        ?= :4000
API_ADDR    ?= :8080
DB_DSN      ?= $(JOURNAL_DB_DSN)
JWT_SECRET  ?= change-me-in-production
CORS_ORIGIN ?= http://localhost:5173

# ------------------------------------
# Helpers
# ------------------------------------
.PHONY: help
help:
	@grep -E '^[a-zA-Z_/.-]+:.*?##' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  %-28s %s\n", $$1, $$2}'

# ====================================
# Code quality
# ====================================

.PHONY: fmt
fmt: ## Format all Go source files
	go fmt ./...

.PHONY: vet
vet: fmt ## Run go vet across all packages
	go vet ./...

# ====================================
# Testing
# ====================================

.PHONY: test
test: vet ## Run all unit tests (skips smoke tests that need TEST_DSN)
	go test ./...

.PHONY: test/v
test/v: vet ## Run all unit tests with verbose output
	go test -v ./...

.PHONY: test/api
test/api: ## Run backend API unit + contract tests
	go test ./backend/...

.PHONY: test/smoke
test/smoke: ## Run smoke tests against a live DB (requires TEST_DSN env var)
	@if [ -z "$(TEST_DSN)" ]; then \
		echo "ERROR: TEST_DSN is not set.  Example:"; \
		echo "  make test/smoke TEST_DSN='postgres://user:pass@localhost/blessed_bites'"; \
		exit 1; \
	fi
	TEST_DSN="$(TEST_DSN)" go test -v -run TestSmoke ./backend/cmd/api/

.PHONY: test/frontend
test/frontend: ## Run SvelteKit type-checks
	cd frontend && npm run check

# ====================================
# Legacy SSR monolith (cmd/web)
# ====================================

.PHONY: run
run: vet ## Run the legacy SSR server locally on $(ADDR)
	go run ./cmd/web -addr="$(ADDR)" -dsn="$(DB_DSN)"

.PHONY: dev
dev: ## Run the legacy SSR server with Air hot-reload
	@if [ -x ./bin/air ]; then \
		APP_ENV="$(APP_ENV)" DB_DSN="$(DB_DSN)" ./bin/air -c .air.toml; \
	else \
		APP_ENV="$(APP_ENV)" DB_DSN="$(DB_DSN)" air -c .air.toml; \
	fi

# ====================================
# New Go REST API  (backend/cmd/api)
# ====================================

.PHONY: run/api
run/api: vet ## Run the REST API server on $(API_ADDR)
	go run ./backend/cmd/api \
		-addr="$(API_ADDR)" \
		-dsn="$(DB_DSN)" \
		-jwt-secret="$(JWT_SECRET)" \
		-cors-origin="$(CORS_ORIGIN)"

.PHONY: build/api
build/api: vet ## Build the REST API binary to bin/api
	go build -o bin/api ./backend/cmd/api

# ====================================
# SvelteKit frontend
# ====================================

.PHONY: frontend/install
frontend/install: ## Install frontend npm dependencies
	cd frontend && npm install

.PHONY: frontend/dev
frontend/dev: ## Start the SvelteKit dev server (hot-reload, port 5173)
	cd frontend && npm run dev

.PHONY: frontend/build
frontend/build: ## Build the SvelteKit app for production
	cd frontend && npm run build

# ====================================
# Docker / full-stack
# ====================================

.PHONY: docker/up
docker/up: ## Start the full refactored stack (API + frontend + Caddy + Postgres)
	docker compose -f docker-compose.refactor.yml up --build

.PHONY: docker/down
docker/down: ## Tear down the refactored stack
	docker compose -f docker-compose.refactor.yml down

.PHONY: docker/up/legacy
docker/up/legacy: ## Start the original monolith compose stack
	docker compose up --build

.PHONY: docker/down/legacy
docker/down/legacy: ## Tear down the original monolith compose stack
	docker compose down

# ====================================
# Database migrations
# ====================================

.PHONY: db/migrations/new
db/migrations/new: ## Create a new migration: make db/migrations/new name=<name>
	migrate create -seq -ext=.sql -dir ./migrations $(name)

.PHONY: db/migrations/up
db/migrations/up: ## Apply all pending migrations
	migrate -path ./migrations -database "$(DB_DSN)" up

.PHONY: db/migrations/down
db/migrations/down: ## Roll back the last migration
	migrate -path ./migrations -database "$(DB_DSN)" down 1
