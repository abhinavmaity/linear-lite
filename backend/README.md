# Backend Runtime Quickstart

This folder contains the Linear-lite backend runtime skeleton for Milestone 1.

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

