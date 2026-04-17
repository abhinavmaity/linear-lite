# Linear-lite

Linear-lite is a focused issue tracking and project planning application for small engineering teams.

It provides the core workflows teams need every day without the complexity of heavyweight enterprise tools.

## Core Capabilities

- Authentication (register, login, session restore)
- Issue lifecycle management (create, update, archive, restore)
- List and board views backed by the same issue source of truth
- Project, sprint, and label management
- Dashboard metrics and recent activity
- Search, filtering, sorting, and pagination for issues

## Tech Stack

- Frontend: React, TypeScript, Vite, Zustand, TanStack Query
- Backend: Go, Gin, PostgreSQL, Redis
- Local runtime: Docker Compose

## Quick Start

From repository root:

```bash
docker compose up --build -d
```

Application endpoints:

- Frontend: `http://localhost:5173`
- Backend API: `http://localhost:8080/api/v1`

If ports are busy, override host ports:

```bash
FULLSTACK_FRONTEND_PORT=5180 FULLSTACK_BACKEND_PORT=18080 docker compose up --build -d
```

`docker compose up` now runs schema migrations automatically before backend startup.

Stop local stack:

```bash
docker compose down -v
```

## Deploy on Railway

Create one Railway project with four services:

- `backend` (root directory: `/backend`)
- `frontend` (root directory: `/frontend`)
- `Postgres` (Railway template)
- `Redis` (Railway template)

Set each app service to use its config-as-code file:

- backend config path: `/backend/railway.json`
- frontend config path: `/frontend/railway.json`

Backend service variables:

- `APP_ENV=production`
- `JWT_SECRET=<at least 32 characters>`
- `DATABASE_URL=${{Postgres.DATABASE_URL}}`
- `REDIS_URL=${{Redis.REDIS_URL}}`
- `MIGRATIONS_PATH=/app/migrations`
- `MIGRATION_DIRECTION=up`
- `MIGRATION_STEPS=0`
- `CORS_ORIGINS=https://<frontend-public-domain>`

Frontend service variable:

- `VITE_API_BASE_URL=https://<backend-public-domain>/api/v1`

Notes:

- Backend migrations are configured in `/backend/railway.json` via pre-deploy command `migrate`.
- Frontend now builds static assets and serves them with `serve` (production runtime, no Vite dev server).

## Local Development

Frontend (`frontend/`):

```bash
npm install
npm run dev
npm run build
```

Backend (`backend/`):

```bash
go mod tidy
go build ./...
go run ./cmd/api
```

## Validation

Common validation commands:

```bash
./scripts/smoke_issue_workflow.sh
./scripts/smoke_cache.sh
```

CI validation is defined in:

- `.github/workflows/ci-validation.yml`

## Repository Structure

- `frontend/` - React application
- `backend/` - Go API and migrations
- `docs/` - architecture and product/delivery documentation
- `scripts/` - smoke validation scripts

## Documentation

- [Technical_Architecture.md](/Users/abhinavmaity/code/linear-lite/docs/Technical_Architecture.md) - backend and API source of truth
- [Product_and_Delivery.md](/Users/abhinavmaity/code/linear-lite/docs/Product_and_Delivery.md) - product scope, UX journeys, delivery summary, and readiness
