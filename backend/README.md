# Backend Runtime Quickstart

This folder contains the Linear-lite backend runtime through Milestone 2.

Implemented in backend so far:
- Milestone 1: runtime foundation (config, server bootstrap, middleware, migration runner)
- Milestone 2: auth database foundation and core auth endpoints
  - `POST /api/v1/auth/register`
  - `POST /api/v1/auth/login`
  - `GET /api/v1/auth/me`

## Environment Variables

Required:
- `DATABASE_URL`
- `REDIS_URL`
- `JWT_SECRET`
- `CORS_ORIGINS`

Optional with defaults:
- `APP_ENV=development`
- `PORT=8080`
- `JWT_TTL=24h`
- `LOG_LEVEL=info`
- `BCRYPT_COST=12`
- `HTTP_READ_HEADER_TIMEOUT=10s`
- `HTTP_READ_TIMEOUT=15s`
- `HTTP_WRITE_TIMEOUT=30s`
- `HTTP_IDLE_TIMEOUT=60s`
- `HTTP_SHUTDOWN_TIMEOUT=10s`
- `DB_MAX_OPEN_CONNS=25`
- `DB_MAX_IDLE_CONNS=10`
- `DB_CONN_MAX_LIFETIME=30m`
- `DB_CONN_MAX_IDLE_TIME=5m`

Migration command variables:
- `MIGRATIONS_PATH=migrations` (default)
- `MIGRATION_DIRECTION=up` (`up` or `down`)
- `MIGRATION_STEPS=0` (`0` means all available steps)

## Commands

From `backend/`:

```bash
go mod tidy
go build ./...
go run ./cmd/api
```

Run migrations explicitly:

```bash
go run ./cmd/migrate
```

Run a limited rollback:

```bash
MIGRATION_DIRECTION=down MIGRATION_STEPS=1 go run ./cmd/migrate
```

## Auth Endpoints (Milestone 2)

Public:
- `POST /api/v1/auth/register`
- `POST /api/v1/auth/login`

Protected:
- `GET /api/v1/auth/me` with `Authorization: Bearer <token>`

Behavior:
- Access token only (no refresh token, no cookie auth)
- Passwords hashed with bcrypt
- Email uniqueness enforced case-insensitively
- Error codes follow architecture envelopes (`validation_error`, `unauthorized`, `conflict`, `internal_error`)

## Response Envelopes

- Success (resource): `{ "data": ... }`
- Success (collection): `{ "items": [...], "pagination": ... }`
- Error: `{ "error": { "code": "...", "message": "...", "fields"?: {...}, "request_id"?: "..." } }`

Shared helpers:
- `internal/errors` centralizes `AppError` and HTTP code mapping.
- `internal/handlers/response.go` provides success envelope writers.
- `internal/validation/helpers.go` provides reusable UUID/pagination/sort/date/repeated-query helpers.

## Docker Baseline

From repo root:

```bash
docker compose -f docker-compose.backend.yml up --build
```

Then run migrations in a one-off backend container:

```bash
docker compose -f docker-compose.backend.yml run --rm backend migrate
```

## Manual Auth Validation Checklist

After startup and migrations, verify this sequence:
1. `POST /api/v1/auth/register` -> `201`
2. Duplicate `POST /api/v1/auth/register` with same email -> `409`
3. `POST /api/v1/auth/login` with valid credentials -> `200`
4. `POST /api/v1/auth/login` with invalid credentials -> `401`
5. `GET /api/v1/auth/me` without token -> `401`
6. `GET /api/v1/auth/me` with bearer token from login -> `200`
