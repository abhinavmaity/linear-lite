# Milestone 6 Validation Report

Date: April 12, 2026  
Scope: M6-03 baseline stability run for current validation commands

## Summary

Baseline command execution for Milestone 6 is now passing end-to-end. Compile/build checks passed, issue workflow smoke passed, and cache smoke passed after fixing a smoke-data collision in the cache script.

## Commands Executed

1. Frontend build
- Command: `npm run build` (cwd: `frontend/`)
- Result: PASS

2. Backend build
- Command: `GOCACHE=/tmp/go-build go build ./...` (cwd: `backend/`)
- Result: PASS

3. Issue workflow smoke
- Command: `./scripts/smoke_issue_workflow.sh` (cwd: repo root)
- Result: PASS

4. Cache smoke
- Command: `./scripts/smoke_cache.sh` (cwd: repo root)
- Result: PASS

## Task Status Impact

- M6-01: Complete
- M6-02: Complete
- M6-03: Complete
- M6-04: Complete
- M6-05: Complete
- M6-06: Complete
- M6-07: Complete
- M6-08: Complete
- M6-09: Complete
- M6-10: Complete
- M6-11: Complete

## Unblock Step

## Investigation Outcome

Root cause:
- `scripts/smoke_cache.sh` updated labels to a fixed name (`cache-label-updated`).
- Labels enforce global case-insensitive uniqueness (`uq_labels_lower_name`).
- On reused database state, this could produce legitimate `409 conflict` responses, causing a flaky smoke assertion.

Fix applied:
- Updated `scripts/smoke_cache.sh` to use a unique update target name per run (`cache-label-updated-<timestamp>`).

Post-fix rerun:

```bash
./scripts/smoke_issue_workflow.sh
./scripts/smoke_cache.sh
```

Both pass after the fix.

## Next-Set Validation (M6-04 to M6-06)

Added tests:
- `backend/internal/handlers/handlers_test.go`
- `backend/internal/services/services_rules_test.go`
- `backend/internal/repositories/issue_repository_test.go`

Command executed:

```bash
GOCACHE=/tmp/go-build go test ./...
```

Result:
- PASS (handlers/services/repositories packages with new tests)

## Next-Set Validation (M6-07 to M6-11)

Harness and tests added:
- `frontend/playwright.config.ts`
- `frontend/tests/e2e/helpers.ts`
- `frontend/tests/e2e/auth-session.spec.ts`
- `frontend/tests/e2e/issue-workflow.spec.ts`
- `frontend/tests/e2e/supporting-resources.spec.ts`
- `frontend/tests/e2e/dashboard-consistency.spec.ts`

Frontend package updates:
- `frontend/package.json` scripts for E2E execution
- `@playwright/test` added to dev dependencies

Command executed:

```bash
npm run e2e
```

Result:
- PASS (5/5 Playwright specs)
  - `M6-08` auth/session journey: pass
  - `M6-09` core issue workflow journey: pass
  - `M6-10` supporting resources journey: pass
  - `M6-11` dashboard consistency journey: pass

Notes:
- E2E server is configured with explicit API base override:
  - `VITE_API_BASE_URL=http://127.0.0.1:8080/api/v1`
- E2E base URL uses `http://localhost:5173` to align with backend CORS settings.

## Next-Set Validation (M6-12 to M6-14)

Infrastructure/docs added:
- `docker-compose.fullstack.yml`
- `frontend/Dockerfile`
- `frontend/.dockerignore`
- `frontend/.env.example`
- README updates for full-stack run + env contract
- backend README updates for full-stack run + migration path

Commands executed:

```bash
FULLSTACK_BACKEND_PORT=18080 FULLSTACK_FRONTEND_PORT=5180 docker compose -f docker-compose.fullstack.yml up -d --build
FULLSTACK_BACKEND_PORT=18080 FULLSTACK_FRONTEND_PORT=5180 docker compose -f docker-compose.fullstack.yml --profile tools run --rm migrate
FULLSTACK_BACKEND_PORT=18080 FULLSTACK_FRONTEND_PORT=5180 docker compose -f docker-compose.fullstack.yml up -d backend frontend
curl http://localhost:5180
curl http://localhost:18080/api/v1/auth/me
```

Results:
- PASS: full stack starts with postgres/redis/backend/frontend containers healthy.
- PASS: migration runner completes via explicit one-off `migrate` service.
- PASS: frontend endpoint returns `200`.
- PASS: backend auth-me endpoint returns `401 unauthorized` envelope when unauthenticated.

Notes:
- Host port collisions existed for default `5173`/`8080` in this environment; validated override path (`5180`/`18080`) and documented it as the recommended fallback.
- Postgres/Redis are intentionally internal-only in full-stack compose (no host port exposure required for MVP local runtime).
