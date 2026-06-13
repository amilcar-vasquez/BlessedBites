# BlessedBites

BlessedBites is a full-stack restaurant management application with a Go REST API backend and a modern Svelte frontend. It supports public browsing and checkout flows, plus authenticated admin tooling for managing users, categories, and menu items.

This README is intended to document how the application is organized so new contributors can quickly understand the codebase and development workflow.

## What The Site Does

- Public pages: home, signup/login, menu browsing, search, category views, checkout.
- Admin pages: dashboard, user management, menu/category management.
- Ordering flow: create orders and order items from checkout.
- Customer feedback: ratings endpoints for menu items.
- Account support: password reset request + reset completion.

## Tech Stack

**Backend:**
- Language: Go 1.24+
- Framework: `net/http` with custom routing
- Database: PostgreSQL via `github.com/lib/pq`
- Sessions: `github.com/gorilla/sessions` cookie store
- API Contract: OpenAPI 3.0 (see `backend/openapi.yaml`)
- Dev tooling: `make`, `migrate`

**Frontend:**
- Framework: SvelteKit
- Build tool: Vite
- Language: TypeScript
- Styling: CSS (responsive design)
- Package manager: npm

**Infrastructure:**
- Reverse proxy: Caddy
- Container runtime: Docker & Docker Compose
- Database: PostgreSQL 16

## High-Level Architecture

The application follows a separated backend/frontend architecture:

**Backend API** (`backend/`):
1. `backend/cmd/api/main.go` initializes the API server, database connection, and middleware.
2. `backend/internal/handlers/` contains HTTP request handlers organized by feature.
3. `backend/internal/data/` contains database models and query logic.
4. `backend/internal/middleware/` provides CORS, auth, logging, and other middleware.
5. `backend/openapi.yaml` defines the REST API contract.

**Frontend** (`frontend/`):
1. `frontend/src/routes/` contains SvelteKit route components (file-based routing).
2. `frontend/src/lib/` contains reusable components and utilities.
3. Frontend communicates with backend via REST API at `/api/v1` (proxied by Caddy).
4. Built with Vite and deployed as a static site preview server in Docker.

**Reverse Proxy** (Caddy):
- Routes requests to frontend on `/` and backend API on `/api/v1`.
- Serves uploaded images from `/uploads/*`.

## Project Structure

```text
.
|- backend/
|  |- cmd/api/
|  |  \- main.go              # API server bootstrap and setup
|  |- internal/
|  |  |- handlers/            # HTTP request handlers by feature
|  |  |- middleware/          # CORS, auth, logging, realtime middleware
|  |  |- data/                # DB models (users, menu, orders, analytics, ratings)
|  |  |- realtime/            # WebSocket and real-time features
|  |  \- utils/              # Utility helpers
|  |- pkg/
|  |  |- jwt/                 # JWT authentication helpers
|  |  \- responses/          # API response formatting
|  |- contract/
|  |  \- openapi_contract_test.go  # API contract validation tests
|  |- openapi.yaml            # OpenAPI 3.0 specification
|  \- Dockerfile              # Backend container build
|-
|- frontend/
|  |- src/
|  |  |- routes/              # SvelteKit file-based routes
|  |  |- lib/                 # Reusable components and utilities
|  |  |- app.html             # HTML entry point
|  |  \- service-worker.ts    # PWA service worker
|  |- static/
|  |  |- manifest.webmanifest # PWA manifest
|  |  \- icons/              # App icons
|  |- package.json            # Frontend dependencies
|  |- svelte.config.js        # SvelteKit configuration
|  |- vite.config.ts          # Vite build configuration
|  \- Dockerfile              # Frontend container build
|-
|- migrations/                # SQL migrations (up/down)
|- deploy/
|  \- Caddyfile              # Caddy reverse proxy configuration
|- documentation/             # SQL dumps and reference artifacts
|- docker-compose.yml         # Original monolithic app orchestration
|- docker-compose.refactor.yml # New separated backend/frontend orchestration
|- Makefile                   # Local/dev/test/migration commands
\- test_csrf.go               # Standalone CSRF-focused test script
```

## API Routes

The backend API is documented in `backend/openapi.yaml` (OpenAPI 3.0). All API endpoints are prefixed with `/api/v1`.

**Common API endpoints** (all routes are RESTful):

- Users: `GET|POST /users`, `GET|PUT /users/{id}`, `DELETE /users/{id}`
- Authentication: `POST /auth/login`, `POST /auth/logout`, `POST /auth/register`
- Menu Items: `GET|POST /menu-items`, `GET|PUT|DELETE /menu-items/{id}`
- Categories: `GET|POST /categories`, `GET|PUT|DELETE /categories/{id}`
- Orders: `GET|POST /orders`, `GET|PUT /orders/{id}`
- Ratings: `GET|POST /ratings`, `GET /ratings/item/{item_id}`
- Analytics: `POST /analytics/events`, `GET /analytics/summary`

For the complete API specification, see `backend/openapi.yaml`.

## Frontend Routes

The frontend uses SvelteKit's file-based routing in `frontend/src/routes/`.

**Public pages:**
- `/` - Home
- `/signup` - Sign up form
- `/login` - Login form
- `/menu` - Menu browsing
- `/checkout` - Shopping cart and checkout
- `/reset-password` - Password reset flow

**Admin pages** (authenticated):
- `/dashboard` - Admin dashboard
- `/menu/manage` - Menu management
- `/users` - User management

## Reverse Proxy Routes (Caddy)

Caddy (`deploy/Caddyfile`) routes:
- `/` → Frontend (port 3000)
- `/api/v1/*` → Backend API (port 8080)
- `/uploads/*` → Static uploaded images

## Backend Data Layer

`backend/internal/data/` contains one model file per domain:

- `user.go`: user CRUD and auth-related queries.
- `category.go`: category CRUD operations.
- `menuItem.go`: menu item CRUD and filtering/search support.
- `order.go`: order creation and retrieval.
- `orderItem.go`: order line items persistence.
- `rating.go`: ratings create/read and aggregation.
- `analytics.go`: event/metrics persistence and queries.
- `recommendations.go`: recommendation data support.

## Frontend Components

`frontend/src/lib/` and `frontend/src/routes/` contain:

- Page-level components in `routes/` for each public/admin route.
- Reusable UI components in `lib/components/`.
- API client utilities in `lib/api/` for communicating with backend.
- Stores in `lib/stores/` for shared state (auth, cart, etc.).
- Utility functions in `lib/utils/`.

## Frontend Assets

**Svelte Components** (`frontend/src/`):

- Route components for each page in `routes/`.
- Shared layout in `+layout.svelte`.
- Reusable UI components in `lib/components/`.

**Static Assets** (`frontend/static/`):

- PWA manifest: `manifest.webmanifest`.
- App icons: `icons/`.

**Uploaded Images** (`uploads/`):

- User-uploaded menu/category images (mounted as Docker volume).
- Served at `/uploads/*` via Caddy.

## Security

**Backend API Authentication:**
- JWT-based authentication via `backend/pkg/jwt/`.
- Token provided in `Authorization: Bearer <token>` header.
- Admin routes protected by JWT validation in `backend/internal/middleware/`.

**CORS:**
- CORS middleware in `backend/internal/middleware/` allows frontend requests.
- Configured to accept frontend origin (e.g., `http://localhost` in development).

**Sessions:**
- Stateless JWT tokens (no cookie-based sessions).
- Tokens are short-lived; refresh tokens may be used for longer sessions.

**Frontend:**
- JWT tokens stored securely (localStorage or secure storage).
- Sent with all API requests via `Authorization` header.

**Development:**
- HTTPS not enforced in development; set `APP_ENV=development`.
- In production, HTTPS is enforced via Caddy.

## Configuration

**Backend Environment Variables** (used by `backend/cmd/api`):

- `APP_ENV`: `development` or production (controls logging, HTTPS enforcement).
- `DB_DSN`: PostgreSQL connection string (e.g., `postgres://user:pass@localhost/db?sslmode=disable`).
- `JWT_SECRET`: Secret key for signing JWT tokens.
- `CORS_ORIGIN`: Allowed frontend origin (e.g., `http://localhost` in dev, full URL in prod).
- `MAIL_HOST`, `MAIL_USERNAME`, `MAIL_PASSWORD`: SMTP settings for password reset emails.

**Frontend Environment Variables** (used during build):

- `PUBLIC_API_BASE_URL`: Base URL for API calls (e.g., `/api/v1` when behind Caddy).

**Docker Compose Variables** (in `.env` or environment):

- Used by `docker-compose.refactor.yml` to configure database, API, and frontend services.
- See `docker-compose.refactor.yml` for required variables.

**Backend CLI Flags** (optional):

- `-addr`: Server address (default `:8080`)
- `-dsn`: Database DSN (default from `DB_DSN` env var)

## Local Development

**Prerequisites:**

- Go 1.24+ (for backend)
- Node.js 20+ and npm (for frontend)
- PostgreSQL 13+ (local or via Docker)
- `migrate` CLI (for SQL migrations) - install with `go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest`

**Setup:**

1. Create `.env` in project root:

```env
APP_ENV=development
DB_DSN=postgresql://blessed_bites:blessed_bites@localhost:5432/blessed_bites?sslmode=disable
JWT_SECRET=dev-secret-change-in-production
CORS_ORIGIN=http://localhost:5173
MAIL_HOST=<smtp-host>
MAIL_USERNAME=<smtp-user>
MAIL_PASSWORD=<smtp-password>
```

2. Start PostgreSQL (local or Docker):

```bash
# Using Docker
docker run -d --name blessedbites-postgres \
  -e POSTGRES_DB=blessed_bites \
  -e POSTGRES_USER=blessed_bites \
  -e POSTGRES_PASSWORD=blessed_bites \
  -p 5432:5432 \
  postgres:16
```

3. Run migrations:

```bash
make db/migrations/up
```

4. **Terminal 1: Start the backend with live-reload:**

```bash
make dev
```

Backend will be available at `http://localhost:8080`.

5. **Terminal 2: Start the frontend dev server:**

```bash
cd frontend
npm install
npm run dev
```

Frontend will be available at `http://localhost:5173`.

**Other useful commands:**

- `make fmt` - format backend code.
- `make vet` - run Go static checks.
- `make run/tests` - run backend test suite.
- `make db/migrations/down` - rollback one migration.
- `cd frontend && npm run build` - build frontend for production.
- `cd frontend && npm run preview` - preview production build locally.

## Docker

**Modern Setup** (recommended for full-stack development):

Use `docker-compose.refactor.yml` which includes backend API, frontend, PostgreSQL, and Caddy reverse proxy:

```bash
docker compose -f docker-compose.refactor.yml up --build
```

Access the application at `http://localhost:8080` (via Caddy).

Services:
- **Frontend** (Svelte): http://localhost:3000 (or via Caddy on 8080)
- **Backend API**: http://localhost:8080/api/v1 (or direct at :8080)
- **Database**: PostgreSQL on :5432
- **Caddy**: Reverse proxy on :8080

**Legacy Setup** (original monolithic approach):

`docker-compose.yml` defines a single app service (older configuration):

```bash
docker compose up --build
```

This maps `4000:4000` and is no longer actively used.

## Database And Documentation Assets

- SQL migrations: `migrations/*.sql`
- SQL dumps: `documentation/blessed_bites_dump.sql`, `documentation/blessed_bites_backup.sql`

## Notes For Contributors

**Backend API:**
- Start from `backend/cmd/api/main.go` to understand API initialization.
- Check `backend/openapi.yaml` for the API contract before implementing endpoints.
- Add new database changes through `migrations/` and update models in `backend/internal/data/`.
- Handlers in `backend/internal/handlers/` should follow the OpenAPI spec.
- Use JWT middleware from `backend/internal/middleware/` for authenticated routes.
- Test API contract compliance with `backend/contract/openapi_contract_test.go`.

**Frontend:**
- Start from `frontend/src/routes/` to understand the page structure (SvelteKit file-based routing).
- Reusable components go in `frontend/src/lib/components/`.
- API communication logic goes in `frontend/src/lib/api/`.
- Use `frontend/src/lib/stores/` for shared state (auth token, user info, etc.).
- The `PUBLIC_API_BASE_URL` environment variable controls API endpoint base (default `/api/v1`).

**General:**
- Keep API and frontend changes aligned with the OpenAPI spec.
- Run `docker compose -f docker-compose.refactor.yml up --build` to test the full stack.
- Document new endpoints in `backend/openapi.yaml`.
