# ====================================
# Project: Blessed Bites Go Web App
# ====================================

# Load environment variables from .env if it exists
ifneq (,$(wildcard .env))
	include .env
	export
endif

# Default environment variables (can be overridden by .env)
APP_ENV ?= development
ADDR ?= :4000
DB_DSN ?= $(JOURNAL_DB_DSN)

# ------------------------------------
# Go Commands
# ------------------------------------

.PHONY: fmt
fmt:
	@echo "🧹 Formatting code..."
	go fmt ./...

.PHONY: vet
vet: fmt
	@echo "🔍 Running go vet..."
	go vet ./...

.PHONY: run/tests
run/tests: vet
	@echo "🧪 Running all tests..."
	go test -v ./...

.PHONY: run
run: vet
	@echo "🚀 Starting app locally on $(ADDR)"
	@echo "   Using DB_DSN=$(DB_DSN)"
	go run ./cmd/web -addr="$(ADDR)" -dsn="$(DB_DSN)"

# ------------------------------------
# Development Mode (with Air)
# ------------------------------------
.PHONY: dev
dev:
	@echo "🔥 Running in development mode with Air..."
	@echo "   APP_ENV=$(APP_ENV)"
	@echo "   DB_DSN=$(DB_DSN)"
	@if [ -x ./bin/air ]; then \
		APP_ENV="$(APP_ENV)" DB_DSN="$(DB_DSN)" ./bin/air -c .air.toml; \
	else \
		APP_ENV="$(APP_ENV)" DB_DSN="$(DB_DSN)" air -c .air.toml; \
	fi

.PHONY: dev-local
dev-local:
	@echo "🚀 Running local dev server (inline env, non-air)"
	@echo "   APP_ENV=development"
	@echo "   DB_DSN=$(DB_DSN)"
	@echo "   ADDR=$(ADDR)"
	@echo "--> Note: this runs the real binary with inline env to ensure child process sees APP_ENV."
	@APP_ENV=development DB_DSN="$(DB_DSN)" go run ./cmd/web -addr="$(ADDR)" -dsn="$(DB_DSN)"

# ------------------------------------
# Database Migrations
# ------------------------------------
.PHONY: db/migrations/new
db/migrations/new:
	@echo "🧩 Creating new migration: $(name)"
	migrate create -seq -ext=.sql -dir ./migrations $(name)

.PHONY: db/migrations/up
db/migrations/up:
	@echo "⬆️  Applying all up migrations..."
	migrate -path ./migrations -database "$(DB_DSN)" up

.PHONY: db/migrations/down
db/migrations/down:
	@echo "⬇️  Rolling back last migration..."
	migrate -path ./migrations -database "$(DB_DSN)" down 1
