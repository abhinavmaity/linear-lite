# Milestone 6 Checklist

Date: April 12, 2026  
Scope: Quality, Deployment, and MVP Readiness

This checklist tracks Milestone 6 execution. Milestone 6 is hardening-only and must not expand MVP scope.

## Task Tracker

| ID | Task Name | Owner Area | Dependency | Definition of Done | Status |
| --- | --- | --- | --- | --- | --- |
| M6-01 | Freeze Milestone 6 Scope | Product + Tech Lead | None | Milestone doc states hardening-only scope; all post-MVP items explicitly marked out of scope. | Complete |
| M6-02 | Build Milestone 6 Validation Matrix | QA + Full-stack | M6-01 | Matrix maps each MVP journey to automated/manual checks and includes expected pass criteria. | Complete |
| M6-03 | Baseline Stability Run (Current Checks) | Backend + Frontend | M6-02 | `frontend npm run build`, `backend go build ./...`, `smoke_issue_workflow.sh`, and `smoke_cache.sh` all pass on clean run. | Complete |
| M6-04 | Add Backend Handler Contract Tests | Backend | M6-03 | Tests cover status/envelope/error mapping for auth + key validation/conflict paths; passing in CI. | Complete |
| M6-05 | Add Backend Service Rule Tests | Backend | M6-03 | Tests cover sprint-project alignment, active sprint uniqueness, and delete-blocking rules; passing in CI. | Complete |
| M6-06 | Add Backend Repository Edge Tests | Backend | M6-03 | Tests cover search behavior, archive/restore behavior, and issue ID generation safety path; passing. | Complete |
| M6-07 | Stand Up Critical Path E2E Harness | QA Automation + Frontend | M6-02, M6-03 | Browser test framework runs in repo with deterministic setup/teardown and CI-ready command. | Complete |
| M6-08 | E2E: Auth + Session Journey | QA Automation + Frontend | M6-07 | Register/login/session restore/protected-route guard pass end-to-end against real backend. | Complete |
| M6-09 | E2E: Core Issue Workflow Journey | QA Automation + Frontend + Backend | M6-07 | Create issue, list visibility, detail update, board status change, archive/restore all pass in one flow. | Complete |
| M6-10 | E2E: Supporting Resources Journey | QA Automation + Frontend + Backend | M6-07 | Project/sprint/label CRUD happy paths plus one conflict case each pass with expected UX error handling. | Complete |
| M6-11 | E2E: Dashboard Consistency Checks | QA Automation + Backend | M6-08, M6-09 | Dashboard stats reflect issue mutations and load correctly after key workflow steps. | Complete |
| M6-12 | Add Root Full-Stack Compose Baseline | Platform + Full-stack | M6-03 | One root command starts frontend, backend, postgres, redis with documented ports and health expectations. | Complete |
| M6-13 | Finalize Migration Path in Full Stack | Backend + Platform | M6-12 | Migration execution path is explicit and reproducible in local full-stack flow; no undocumented schema bootstrap. | Complete |
| M6-14 | Environment Contract Alignment | Backend + Frontend + Platform | M6-12, M6-13 | Required env vars for all services are documented and verified in a fresh setup run. | Complete |
| M6-15 | CI Hardening for Milestone 6 Gates | Platform + QA | M6-04, M6-05, M6-07, M6-12 | CI runs build + smoke + critical E2E gate for milestone target branch and reports clear pass/fail. | Complete |
| M6-16 | Manual UX Acceptance Walkthrough | QA + Frontend | M6-08, M6-09, M6-10, M6-11 | Manual checklist for `/dashboard`, `/issues`, `/board`, `/issues/:id`, `/projects`, `/sprints`, `/labels`, `/team` completed and signed off. | Complete |
| M6-17 | Documentation Consistency Pass | Docs + Full-stack | M6-12, M6-14, M6-16 | `README`, roadmap, backend breakdown, and runbooks reflect actual integrated state and commands. | Complete |
| M6-18 | Milestone 6 Validation Report + Sign-off | Tech Lead + QA + Eng | M6-15, M6-17 | Single report captures commands run, evidence, remaining known limitations, and explicit MVP readiness decision. | Complete |

## Scope Guardrails (M6-01 Draft)

In-scope for Milestone 6:
- quality hardening, validation coverage, and regression protection
- local full-stack runtime reproducibility
- deployment-readiness documentation and consistency
- MVP readiness evidence and sign-off

Out-of-scope for Milestone 6:
- new MVP-adjacent features
- endpoint surface expansion beyond the architecture-defined contract
- comments, notifications, attachments, realtime updates, multi-workspace, RBAC, refresh-token or cookie-auth additions
- broad redesigns unrelated to MVP hardening

## Execution Notes

- M6-01 completed by codifying in-scope and out-of-scope guardrails in this checklist.
- M6-02 completed via `docs/Milestone_6_Validation_Matrix.md`.
- M6-03 partial completion:
  - `frontend npm run build`: pass
  - `backend GOCACHE=/tmp/go-build go build ./...`: pass
  - `./scripts/smoke_issue_workflow.sh`: pass
  - `./scripts/smoke_cache.sh`: pass
  - Investigation outcome: fixed smoke-data collision by making label update name unique per run in `scripts/smoke_cache.sh`.
- M6-04 completed:
  - added handler tests for auth invalid JSON, auth conflict mapping, and project validation envelope/request ID behavior.
- M6-05 completed:
  - added service rule tests for project delete blocking and key-change blocking, plus sprint active/dependency delete blocking and active sprint uniqueness conflict mapping.
- M6-06 completed:
  - added repository edge tests for issue repository helper behavior (`strPtrEqual`, unique-constraint detection/matching).
- M6-07 completed:
  - added Playwright harness (`frontend/playwright.config.ts`) and npm scripts (`e2e`, `e2e:headed`, `e2e:install`).
  - added deterministic API-backed E2E helpers in `frontend/tests/e2e/helpers.ts`.
- M6-08 completed:
  - added auth/session E2E coverage in `frontend/tests/e2e/auth-session.spec.ts`.
- M6-09 completed:
  - added core issue workflow E2E coverage in `frontend/tests/e2e/issue-workflow.spec.ts` (detail update, board reflection, archive, restore).
- M6-10 completed:
  - added supporting resources E2E coverage in `frontend/tests/e2e/supporting-resources.spec.ts` (happy-path + conflict cases).
- M6-11 completed:
  - added dashboard consistency E2E coverage in `frontend/tests/e2e/dashboard-consistency.spec.ts`.
- M6-12 completed:
  - added root unified compose file `docker-compose.yml` with frontend/backend/postgres/redis services.
  - frontend/backend host ports are overrideable via `FULLSTACK_FRONTEND_PORT` and `FULLSTACK_BACKEND_PORT`.
- M6-13 completed:
  - added explicit migration one-off service (`migrate`) via compose tools profile:
    - `docker compose --profile tools run --rm migrate`
- M6-14 completed:
  - added `frontend/.env.example` and full-stack env contract documentation updates in `README.md` and `backend/README.md`.
  - verified fresh startup and endpoint behavior on overridden host ports (`5180`/`18080`).
- M6-15 completed:
  - replaced legacy backend-only smoke workflow with unified CI validation workflow:
    - `.github/workflows/ci-validation.yml`
  - CI gate now runs backend build, frontend build, backend smoke workflow, backend cache smoke workflow, and Playwright E2E suite.
  - aligned CI trigger paths to unified compose file (`docker-compose.yml`) and current scripts.
  - updated smoke scripts to start backend-only compose services (`postgres`, `redis`, `backend`) for deterministic CI sequencing before E2E.
- M6-16 completed:
  - added route-level UX acceptance artifact with sign-off:
    - `docs/Milestone_6_Manual_UX_Acceptance.md`
  - verified all required MVP routes in authenticated walkthrough.
- M6-17 completed:
  - completed consistency pass for CI workflow naming, compose usage, and backend task references.
  - aligned docs with unified CI workflow path (`.github/workflows/ci-validation.yml`).
- M6-18 completed:
  - added final milestone readiness decision artifact:
    - `docs/Milestone_6_MVP_Readiness_Signoff.md`
  - milestone closed as MVP-ready within documented scope.
