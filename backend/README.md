# BlessedBites API (Transition Scaffold)

This folder contains the API-first transition scaffold.

## Run locally

```bash
go run ./backend/cmd/api -addr=:8080 -dsn="$DB_DSN" -jwt-secret="$JWT_SECRET"
```

## Implemented endpoints

Public:

- `GET /api/v1/health`
- `GET /api/v1/menu`
- `GET /api/v1/menu/{id}`
- `GET /api/v1/categories`
- `GET /api/v1/search?q=`
- `POST /api/v1/orders`
- `GET /api/v1/orders/stream` (SSE)
- `POST /api/v1/ratings`
- `GET /api/v1/ratings/{menu_item_id}`

Auth:

- `POST /api/v1/auth/login`
- `POST /api/v1/auth/signup`
- `POST /api/v1/auth/logout`
- `POST /api/v1/auth/refresh`
- `POST /api/v1/auth/reset-password-request`
- `POST /api/v1/auth/reset-password`

Admin (JWT + role=admin required):

- `GET /api/v1/admin/menu`
- `POST /api/v1/admin/menu`
- `PUT /api/v1/admin/menu/{id}`
- `DELETE /api/v1/admin/menu/{id}`
- `POST /api/v1/admin/category`
- `DELETE /api/v1/admin/category/{id}`
- `GET /api/v1/admin/orders`

## Notes

- Domain logic currently reuses existing models in `internal/data`.
- Existing `cmd/web` remains untouched for safe migration rollout.
- API contract source: `backend/openapi.yaml`.
- Contract guard test: `backend/contract/openapi_contract_test.go`.
