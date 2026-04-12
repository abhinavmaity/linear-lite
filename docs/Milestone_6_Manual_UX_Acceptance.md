# Milestone 6 Manual UX Acceptance

Date: April 13, 2026  
Scope: M6-16 manual walkthrough for `/dashboard`, `/issues`, `/board`, `/issues/:id`, `/projects`, `/sprints`, `/labels`, `/team`

## Environment

- Full stack started with:
  - `FULLSTACK_FRONTEND_PORT=5180 FULLSTACK_BACKEND_PORT=18080 docker compose up -d --build`
- Migrations applied with:
  - `FULLSTACK_FRONTEND_PORT=5180 FULLSTACK_BACKEND_PORT=18080 docker compose --profile tools run --rm migrate`
- Walkthrough base URLs:
  - Frontend: `http://localhost:5180`
  - Backend: `http://localhost:18080/api/v1`

## Walkthrough Method

- Seeded a fresh user + project + sprint + label + issue via API.
- Performed a browser walkthrough (headless Chromium) with authenticated session token in local storage.
- For each route, verified:
  - expected page identity text is present
  - no generic failure content (`Something went wrong`, `Internal Server Error`)
  - URL resolves to the expected route

## Route Checklist

| Route | Expected Marker | Result |
| --- | --- | --- |
| `/dashboard` | `Dashboard` | Pass |
| `/issues` | `Issues` | Pass |
| `/board` | `Board` | Pass |
| `/issues/:id` | `Issue Detail` | Pass |
| `/projects` | `Projects` | Pass |
| `/sprints` | `Sprints` | Pass |
| `/labels` | `Labels` | Pass |
| `/team` | `Team` | Pass |

## Evidence Snapshot

- Seed run id: `1776019215024`
- Issue id used for detail route: `80f92045-35a4-4a53-b535-210a054ef71a`
- All 8 required routes passed in a single authenticated session.

## M6-16 Decision

Manual UX acceptance walkthrough is complete and signed off for Milestone 6 scope.
