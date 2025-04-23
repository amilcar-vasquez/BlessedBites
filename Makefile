## Filename Makefile
include .envrc

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
	go run ./cmd/web -addr=":4000" -dsn=${JOURNAL_DB_DSN}

.PHONY: db/psql
db/psql:
	psql ${JOURNAL_DB_DSN_DB_DSN}

.PHONY: db/migrations/new
db/migrations/new:
		@echo "Creating new migration for $(name)..."
		migrate create -seq -ext=.sql -dir ./migrations $(name)

.PHONY: db/migrations/up
db/migrations/up:
		@echo "Applying all up migrations..."
		migrate -path ./migrations -database ${JOURNAL_DB_DSN} up

.PHONY: dev
dev:
		@echo "Running in development mode..."
		@air