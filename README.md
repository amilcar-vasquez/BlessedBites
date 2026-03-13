# BlessedBites

BlessedBites is a Go web application for restaurant operations. It supports public browsing and checkout flows, plus authenticated admin tooling for managing users, categories, and menu items.

This README is intended to document how the current site is organized so new contributors can quickly understand the codebase.

## What The Site Does

- Public pages: home, signup/login, menu browsing, search, category views, checkout.
- Admin pages: dashboard, user management, menu/category management.
- Ordering flow: create orders and order items from checkout.
- Customer feedback: ratings endpoints for menu items.
- Account support: password reset request + reset completion.
- Server-side rendering: HTML templates from `ui/html`.

## Tech Stack

- Backend: Go (`net/http`) with route patterns from `ServeMux`.
- Database: PostgreSQL via `github.com/lib/pq`.
- Sessions: `github.com/gorilla/sessions` cookie store.
- CSRF: custom middleware in `internal/csrf`.
- Templating: Go `html/template` + shared base template.
- Static assets: CSS, JS, PWA files, uploaded images from `ui/static`.
- Dev tooling: `make`, `air`, `migrate`, Docker.

## High-Level Architecture

1. `cmd/web/main.go` builds the `application` object (models, logger, sessions, mailer, templates).
2. `cmd/web/server.go` wraps the router with CSRF middleware and starts `http.Server`.
3. `cmd/web/routes.go` registers route handlers and mounts static file servers.
4. Handler files in `cmd/web/*Handlers.go` execute use-case logic.
5. Data models in `internal/data/*.go` handle database operations.
6. Templates in `ui/html/*.tmpl` render server responses.

## Project Structure

```text
.
|- cmd/web/
|  |- main.go                 # App bootstrap, DB connect, dependency wiring
|  |- server.go               # HTTP server + CSRF middleware setup
|  |- routes.go               # Route map and static mounts
|  |- middleware.go           # Request logging, auth guard, cache control
|  |- *Handlers.go            # Route handlers grouped by feature
|  |- templates.go            # Template cache construction
|  |- render.go               # Render helpers for templates
|  |- templateData.go         # Per-request template data shaping
|  \- session_helpers.go      # Safe session access helpers
|- internal/
|  |- data/                   # DB models (users, menu, orders, analytics, ratings)
|  |- csrf/csrf.go            # Custom CSRF token + validation middleware
|  |- mailer/mailer.go        # SMTP mail sender abstraction
|  |- utils/                  # Utility helpers (pagination, random phone)
|  \- validator/             # Validation helpers
|- ui/
|  |- html/                   # SSR templates
|  \- static/                 # CSS/JS/images/uploads/PWA files
|- migrations/                # SQL migrations (up/down)
|- documentation/             # SQL dumps and reference artifacts
|- Dockerfile                 # Container build for app
|- docker-compose.yml         # App container orchestration
|- Makefile                   # Local/dev/test/migration commands
\- test_csrf.go               # Standalone CSRF-focused test script
```

## Route Map (Current)

Defined in `cmd/web/routes.go`.

Public endpoints:

- `GET /`
- `GET /signup`
- `POST /signup/new`
- `GET /signup-thanks`
- `GET /login`
- `POST /login`
- `POST /logout`
- `GET /checkout`
- `POST /orders`
- `POST /ratings`
- `GET /ratings/{menu_item_id}`
- `/search`
- `/search.json`
- `/menu/category/`
- `GET /reset-password-request`
- `POST /reset-password-request`
- `GET /reset-password`
- `POST /reset-password`

Authenticated/admin endpoints (`requireAuth`):

- `GET /users`
- `POST /user/update/form`
- `POST /user/update`
- `POST /users/delete`
- `GET /menu`
- `GET /menu/add`
- `POST /menu/add/new`
- `POST /menu/delete`
- `POST /menu/edit`
- `POST /menu/update`
- `POST /menu/active`
- `GET /dashboard`

Category admin endpoints:

- `POST /category/add`
- `POST /category/delete`

Utility/static endpoints:

- `/static/*` and `/ui/static/*` serve assets from `ui/static`.
- `/favicon.ico` serves `ui/static/img/favicon.ico`.
- `GET /debug/csrf` is available only when `APP_ENV=development`.

## Data Layer By File

`internal/data/` contains one model per domain:

- `user.go`: user CRUD and auth-related queries.
- `category.go`: category CRUD.
- `menuItem.go`: menu item CRUD and filtering/search support.
- `order.go`: order creation and retrieval.
- `orderItem.go`: order line items persistence.
- `rating.go`: ratings create/read aggregation.
- `analytics.go`: event/metrics persistence.
- `recommendations.go`: recommendation data support.

## Templates And Frontend Assets

Templates (`ui/html`):

- Shared layout: `base.tmpl`.
- Feature pages: `home.tmpl`, `menu.tmpl`, `checkout.tmpl`, `dashboard.tmpl`, `users.tmpl`, `login.tmpl`, `signup.tmpl`, `signupThanks.tmpl`, `AddMenuItem.tmpl`, `ResetPassword*.tmpl`.

Static assets (`ui/static`):

- CSS: `css/styles.css`, `css/responsiveStyles.css`.
- JS: `js/app.js`, service worker `js/sw.js`.
- PWA: `manifest.json`, `offline.html`.
- Uploaded images: `img/uploads` (mounted as a Docker volume).

## Security And Session Behavior

- Sessions are cookie-based using Gorilla sessions.
- Session keys can be provided with `SESSION_KEYS` (comma-separated, supports key rotation).
- Fallback key source is `SESSION_KEY`.
- CSRF middleware wraps all routes in `cmd/web/server.go`.
- Cookie behavior in development uses `Secure=false`.
- Cookie behavior outside development uses `Secure=true`.
- `requireAuth` middleware protects admin routes.
- In development, `requireAuth` currently auto-authenticates a local admin session.

## Configuration

Environment variables used by the app:

- `APP_ENV`: expected values include `development` and production-like values.
- `DB_DSN`: PostgreSQL DSN consumed by `-dsn` flag default.
- `JOURNAL_DB_DSN`: used by `Makefile` to populate `DB_DSN`.
- `SESSION_KEY` / `SESSION_KEYS`: session and CSRF key sources.
- `MAIL_HOST`, `MAIL_USERNAME`, `MAIL_PASSWORD`: SMTP settings for password reset email.

Flags supported by `cmd/web`:

- `-addr` (default `:4000`)
- `-dsn` (default from `DB_DSN`)

## Local Development

Prerequisites:

- Go toolchain installed.
- PostgreSQL database reachable from your machine.
- `migrate` CLI if running SQL migrations.
- Optional: `air` binary (or bundled `./bin/air`).

1. Create `.env` in project root:

```env
APP_ENV=development
JOURNAL_DB_DSN=postgresql://<user>:<password>@<host>:<port>/<db>?sslmode=require
SESSION_KEY=<random-long-secret>
MAIL_HOST=<smtp-host>
MAIL_USERNAME=<smtp-user>
MAIL_PASSWORD=<smtp-password>
```

2. Run live-reload dev server:

```bash
make dev
```

3. Or run directly without `air`:

```bash
make dev-local
```

Other useful commands:

- `make fmt` - format code.
- `make vet` - run static checks.
- `make run/tests` - run test suite.
- `make run` - run app using `ADDR` and `DB_DSN`.
- `make db/migrations/up` - apply migrations.
- `make db/migrations/down` - rollback one migration.

## Docker

`docker-compose.yml` currently defines a single app service:

- Builds from local `Dockerfile`.
- Loads `.env`.
- Maps `4000:4000`.
- Mounts `./ui/static/img/uploads` into the container for persistent uploads.
- Sets `DB_DSN` from `JOURNAL_DB_DSN`.

Run:

```bash
docker compose up --build
```

## Database And Documentation Assets

- SQL migrations: `migrations/*.sql`
- SQL dumps: `documentation/blessed_bites_dump.sql`, `documentation/blessed_bites_backup.sql`

## Notes For Contributors

- Start from `cmd/web/main.go` and `cmd/web/routes.go` to trace request flow.
- Add new DB changes through `migrations/` and keep model methods in `internal/data` aligned.
- Keep route-level authorization explicit by wrapping admin handlers with `requireAuth`.
- If changing template names, update template cache/render references in `cmd/web/templates.go` and `cmd/web/render.go`.
