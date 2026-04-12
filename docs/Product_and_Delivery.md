# Product and Delivery

## Purpose

This document consolidates product intent, MVP scope, UX journey expectations, delivery outcomes, and readiness status.

`docs/Technical_Architecture.md` remains the source of truth for backend contracts, schema, constraints, and infrastructure behavior.

## Product Objective

Linear-lite is designed to provide a clear, fast, and predictable issue tracking experience for small engineering teams.

The product prioritizes essential planning and execution workflows over feature breadth.

## Target Users

- Small engineering teams
- Teams that need structured issue/project planning
- Teams that prefer a lightweight self-hostable tool over enterprise-heavy project systems

## MVP Scope

Included:

- Authentication
- Issues (create, update, archive, restore)
- List and board issue views
- Projects
- Sprints
- Labels
- Dashboard summary metrics
- Search/filter/sort/pagination for issues

Deferred (out of current scope):

- Comments/discussions
- Attachments/uploads
- Notifications
- Realtime collaboration
- Advanced analytics
- Multi-workspace and RBAC
- Refresh-token/cookie-auth extensions

## UX Journey Baseline

Primary journeys covered by current implementation:

1. Onboarding and access
- Register
- Login
- Session restore
- Protected route enforcement

2. Daily issue execution
- Create issue
- Manage issue state from list and board
- Edit issue detail
- Archive and restore issue

3. Planning flows
- Create and manage projects
- Create and manage sprints
- Create and manage labels

4. Visibility flows
- Dashboard stats and recent activity
- Team directory (read-only in current scope)

## Delivery Summary

Current repository state:

- Frontend is integrated with real backend APIs for MVP flows.
- Backend endpoints for auth, users, projects, sprints, labels, issues, and dashboard are implemented.
- Redis-backed read caching and invalidation are implemented for supported resources.
- Local full-stack runtime is reproducible with Docker Compose.
- CI validates build, smoke, and E2E coverage in a unified workflow.

## Quality and Validation

Automated checks in regular use:

- Frontend build: `npm run build` (`frontend/`)
- Backend build: `go build ./...` (`backend/`)
- Smoke tests:
  - `./scripts/smoke_issue_workflow.sh`
  - `./scripts/smoke_cache.sh`
- Browser E2E suite: `npm run e2e` (`frontend/`)

CI workflow:

- `.github/workflows/ci-validation.yml`

## Local Runtime Contract

First-time local startup from repo root:

```bash
docker compose up --build -d
docker compose --profile tools run --rm migrate
```

Default local endpoints:

- Frontend: `http://localhost:5173`
- Backend API: `http://localhost:8080/api/v1`

Port overrides are supported with:

- `FULLSTACK_FRONTEND_PORT`
- `FULLSTACK_BACKEND_PORT`

## Readiness Status

The project is MVP-ready within the documented scope and deferrals above.

Future work should be prioritized as explicit post-MVP roadmap items rather than reintroducing broad planning/milestone artifacts into operational docs.
