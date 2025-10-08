## Filename Makefile

.PHONY: run/tests
run/tests: vet
	go test -v ./...

.PHONY: fmt
fmt: 
	go fmt ./...

.PHONY: vet
vet: fmt
	go vet ./...

.PHONY: run
run: vet
	go run ./cmd/web -addr=":4000" -dsn=${DB_DSN}

## local psql target removed; use your own connection to Supabase or psql directly

.PHONY: db/migrations/new
db/migrations/new:
		@echo "Creating new migration for $(name)..."
		migrate create -seq -ext=.sql -dir ./migrations $(name)

.PHONY: db/migrations/up
db/migrations/up:
		@echo "Applying all up migrations..."
		migrate -path ./migrations -database ${DB_DSN} up

.PHONY: dev
dev:
		@echo "Running in development mode..."
		@# Load .env into environment (export all variables) and set DB_DSN from JOURNAL_DB_DSN if present
		@set -a; if [ -f "$$(pwd)/.env" ]; then . "$$(pwd)/.env"; fi; set +a; \
		if [ -n "$$JOURNAL_DB_DSN" ]; then export DB_DSN="$$JOURNAL_DB_DSN"; fi; \
		export APP_ENV=development; \
		# Prefer bundled bin/air if available, otherwise use system 'air'
		if [ -x ./bin/air ]; then ./bin/air -c .air.toml; else air -c .air.toml; fi