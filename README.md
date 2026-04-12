# Linear-lite

Linear-lite is a lightweight issue tracking and project management application for small development teams. The goal is to deliver the most important day-to-day workflows teams need without the overhead and complexity of large enterprise tools.

We are building a product that feels familiar to teams who have used tools like Linear, Jira, or GitHub Issues, but with a tighter MVP scope, faster setup, and a cleaner self-hostable experience. The focus is not on building every possible project management feature. The focus is on building the right core workflows really well.

## What We Are Building

Linear-lite is an MVP issue tracker with:

- user registration and login
- issue creation, editing, assignment, and archiving
- project and sprint organization
- labels and priorities
- list and board views for issues
- filtering and full-text search
- issue activity history
- dashboard-level summary metrics

The intended experience is fast, simple, and predictable. A user should be able to sign up, create a project, create issues, assign work, move issues across statuses, and understand progress without needing a complicated setup flow.

## Aim

The aim of this project is to create a streamlined planning and execution tool for small engineering teams that need structure, but do not want the operational and UX weight of a large-scale project management platform.

At a product level, that means:

- reducing friction in daily issue management
- making sprint and project planning easier to understand
- giving teams multiple views of the same work
- keeping the feature set intentionally focused
- making the system easy to run locally or self-host

## Objectives

- Build a clear MVP with strong implementation boundaries.
- Support the most common engineering team workflows end to end.
- Keep the backend architecture explicit enough that implementation has no ambiguity.
- Keep the frontend screen flows aligned with real user journeys.
- Deliver a system that is simple to maintain, extend, and deploy.

## Result We Are Trying To Achieve

The result we are working toward is a complete, implementation-ready MVP where:

- a team can authenticate and start using the product immediately
- issues can be created, updated, searched, filtered, and organized reliably
- projects and sprints provide planning structure
- board and list views reflect the same source of truth
- activity history gives visibility into issue changes
- the application is backed by a well-defined API and database contract
- the full system can be developed and deployed with confidence because the architecture is already documented in detail

In short, we are trying to ship a focused, high-clarity issue tracker that is practical for real team use and straightforward for engineers to implement.

## MVP Scope

Included in MVP:

- authentication
- issue management
- labels
- projects
- sprints
- dashboard
- list view
- board view
- filtering and search
- Docker-based deployment readiness

Explicitly out of scope for MVP:

- comments and discussions
- file uploads and attachments
- email notifications
- realtime collaboration
- time tracking
- issue dependencies
- advanced analytics
- bulk operations
- external integrations
- mobile apps
- multi-workspace support

## Product Principles

- Essential features only
- Fast and lightweight
- Familiar issue tracking patterns
- Self-hostable and implementation-friendly

## Current Planning Sources

The main planning and architecture references for this repository are:

- [Objective.md](/Users/abhinavmaity/code/linear-lite/docs/Objective.md)
- [Frontend_Planning.md](/Users/abhinavmaity/code/linear-lite/docs/Frontend_Planning.md)
- [Technical_Architecture.md](/Users/abhinavmaity/code/linear-lite/docs/Technical_Architecture.md)
- [Integration_Roadmap.md](/Users/abhinavmaity/code/linear-lite/docs/Integration_Roadmap.md)
- [Backend_Task_Breakdown.md](/Users/abhinavmaity/code/linear-lite/docs/Backend_Task_Breakdown.md)
- [Milestone_5_Parity_Completion_Report.md](/Users/abhinavmaity/code/linear-lite/docs/Milestone_5_Parity_Completion_Report.md)
- [Milestone_5_Validation_Report.md](/Users/abhinavmaity/code/linear-lite/docs/Milestone_5_Validation_Report.md)

## Current Status

The repository has moved beyond pure planning. The frontend core MVP flows are implemented, including auth screens, dashboard, issues list, board, issue detail, a create issue modal, and integrated supporting pages for projects, sprints, labels, and team views.

Backend Milestones 1 through 4 are now implemented:
- Milestone 1: runtime foundation
- Milestone 2: database auth foundation + core auth flow
- Milestone 3: core issue workflow backend
- Milestone 4: dashboard and supporting resource APIs

Milestone 5 (Frontend Integration Parity Pass) is now implemented:
- mock fallback routing removed for active frontend flows
- dashboard, issues list, board, issue detail, and create issue modal parity completed against real API contracts
- projects, sprints, and labels pages moved to real CRUD UI flows with backend conflict/validation handling
- team page kept read-only with improved filtering/sorting and loading/error states
- skeleton-loading integration added with `boneyard-js` for major loading routes

The backend now includes:
- canonical SQL schema support for `users`, `projects`, `sprints`, `labels`, `issues`, `issue_labels`, and `issue_activities`
- auth endpoints: `POST /api/v1/auth/register`, `POST /api/v1/auth/login`, `GET /api/v1/auth/me`
- user endpoints: `GET /api/v1/users`, `GET /api/v1/users/:id`
- project endpoints: `GET /api/v1/projects`, `POST /api/v1/projects`, `GET /api/v1/projects/:id`, `PUT /api/v1/projects/:id`, `DELETE /api/v1/projects/:id`
- sprint endpoints: `GET /api/v1/sprints`, `POST /api/v1/sprints`, `GET /api/v1/sprints/:id`, `PUT /api/v1/sprints/:id`, `DELETE /api/v1/sprints/:id`
- label endpoints: `GET /api/v1/labels`, `POST /api/v1/labels`, `GET /api/v1/labels/:id`, `PUT /api/v1/labels/:id`, `DELETE /api/v1/labels/:id`
- issue workflow endpoints: `GET /api/v1/issues`, `POST /api/v1/issues`, `GET /api/v1/issues/:id`, `PUT /api/v1/issues/:id`, `DELETE /api/v1/issues/:id`
- dashboard endpoint: `GET /api/v1/dashboard/stats`
- Redis-backed cache paths and invalidation for users, projects, sprints, labels, and dashboard read endpoints

Frontend auth flows are wired to the real backend contract (not mock auth): register, login, session restore on refresh, and logout redirect behavior.

Dashboard and supporting resource CRUD backend domains are no longer in-progress and now match the Milestone 4 contract surface.

Frontend loading UIs now support auto-generated skeleton flows via `boneyard-js`:
- dependency installed in `frontend/`
- skeleton wrappers added to key routes
- capture scripts added: `npm run boneyard:build` and `npm run boneyard:watch`
- baseline config: `frontend/boneyard.config.json`

## Implementation Snapshot

- Product definition and architecture: complete
- Frontend core shell and issue workflows: complete and backend-backed
- Frontend auth integration with real backend: complete
- Backend Milestone 1 runtime foundation: complete
- Backend Milestone 2 auth foundation and core auth endpoints: complete
- Backend Milestone 3 core issue workflow backend: complete
- Backend Milestone 4 dashboard and supporting resource APIs: complete
- Frontend Milestone 5 integration parity pass: complete
- Supporting resource screens (projects/sprints/labels): CRUD parity integrated
- Team page: read-only parity integrated
- Skeleton-loading integration: complete for major loading routes
- Deployment hardening and broader QA: pending (Milestone 6)

## Backend Smoke Validation

A reproducible backend issue-workflow smoke script is available at:
- [smoke_issue_workflow.sh](/Users/abhinavmaity/code/linear-lite/scripts/smoke_issue_workflow.sh)
- [smoke_cache.sh](/Users/abhinavmaity/code/linear-lite/scripts/smoke_cache.sh)

Run from repo root:

```bash
./scripts/smoke_issue_workflow.sh
./scripts/smoke_cache.sh
```

## Full-Stack Local Runtime (Milestone 6)

Run from repo root:

```bash
docker compose up --build -d
```

First-time setup (fresh clone / fresh database):

```bash
docker compose up --build -d
docker compose --profile tools run --rm migrate
```

After the first migration run, normal day-to-day startup is just:

```bash
docker compose up --build -d
```

If host ports `5173` or `8080` are busy, override them:

```bash
FULLSTACK_FRONTEND_PORT=5180 FULLSTACK_BACKEND_PORT=18080 docker compose up --build -d
```

Apply migrations (explicit one-off path):

```bash
docker compose --profile tools run --rm migrate
```

If you started compose with overridden ports, reuse the same env vars for migration:

```bash
FULLSTACK_FRONTEND_PORT=5180 FULLSTACK_BACKEND_PORT=18080 docker compose --profile tools run --rm migrate
```

Service URLs after startup:
- Frontend: `http://localhost:5173`
- Backend API: `http://localhost:8080/api/v1`
- Postgres: internal-only (`postgres:5432` in compose network)
- Redis: internal-only (`redis:6379` in compose network)

When using port overrides, frontend/backend URLs use the overridden host ports.

Stop full-stack runtime:

```bash
docker compose down -v
```

### Environment Contract

- Backend env baseline: `backend/.env.example`
- Frontend env baseline: `frontend/.env.example`
- Full-stack compose injects:
  - `VITE_API_BASE_URL=http://localhost:${FULLSTACK_BACKEND_PORT:-8080}/api/v1`
  - `DATABASE_URL=postgres://postgres:postgres@postgres:5432/linear_lite?sslmode=disable`
  - `REDIS_URL=redis://redis:6379/0`
  - `CORS_ORIGINS=http://localhost:${FULLSTACK_FRONTEND_PORT:-5173},http://localhost:3000`

## Frontend Skeleton Capture (Milestone 5)

Run from `frontend/`:

```bash
npm run dev
```

In a second terminal:

```bash
npm run boneyard:build
# or
npm run boneyard:watch
```

This generates responsive bones in `frontend/src/bones/` for routes wrapped with `Skeleton`.
