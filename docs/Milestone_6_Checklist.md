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
| M6-04 | Add Backend Handler Contract Tests | Backend | M6-03 | Tests cover status/envelope/error mapping for auth + key validation/conflict paths; passing in CI. | Not Started |
| M6-05 | Add Backend Service Rule Tests | Backend | M6-03 | Tests cover sprint-project alignment, active sprint uniqueness, and delete-blocking rules; passing in CI. | Not Started |
| M6-06 | Add Backend Repository Edge Tests | Backend | M6-03 | Tests cover search behavior, archive/restore behavior, and issue ID generation safety path; passing. | Not Started |
| M6-07 | Stand Up Critical Path E2E Harness | QA Automation + Frontend | M6-02, M6-03 | Browser test framework runs in repo with deterministic setup/teardown and CI-ready command. | Not Started |
| M6-08 | E2E: Auth + Session Journey | QA Automation + Frontend | M6-07 | Register/login/session restore/protected-route guard pass end-to-end against real backend. | Not Started |
| M6-09 | E2E: Core Issue Workflow Journey | QA Automation + Frontend + Backend | M6-07 | Create issue, list visibility, detail update, board status change, archive/restore all pass in one flow. | Not Started |
| M6-10 | E2E: Supporting Resources Journey | QA Automation + Frontend + Backend | M6-07 | Project/sprint/label CRUD happy paths plus one conflict case each pass with expected UX error handling. | Not Started |
| M6-11 | E2E: Dashboard Consistency Checks | QA Automation + Backend | M6-08, M6-09 | Dashboard stats reflect issue mutations and load correctly after key workflow steps. | Not Started |
| M6-12 | Add Root Full-Stack Compose Baseline | Platform + Full-stack | M6-03 | One root command starts frontend, backend, postgres, redis with documented ports and health expectations. | Not Started |
| M6-13 | Finalize Migration Path in Full Stack | Backend + Platform | M6-12 | Migration execution path is explicit and reproducible in local full-stack flow; no undocumented schema bootstrap. | Not Started |
| M6-14 | Environment Contract Alignment | Backend + Frontend + Platform | M6-12, M6-13 | Required env vars for all services are documented and verified in a fresh setup run. | Not Started |
| M6-15 | CI Hardening for Milestone 6 Gates | Platform + QA | M6-04, M6-05, M6-07, M6-12 | CI runs build + smoke + critical E2E gate for milestone target branch and reports clear pass/fail. | Not Started |
| M6-16 | Manual UX Acceptance Walkthrough | QA + Frontend | M6-08, M6-09, M6-10, M6-11 | Manual checklist for `/dashboard`, `/issues`, `/board`, `/issues/:id`, `/projects`, `/sprints`, `/labels`, `/team` completed and signed off. | Not Started |
| M6-17 | Documentation Consistency Pass | Docs + Full-stack | M6-12, M6-14, M6-16 | `README`, roadmap, backend breakdown, and runbooks reflect actual integrated state and commands. | Not Started |
| M6-18 | Milestone 6 Validation Report + Sign-off | Tech Lead + QA + Eng | M6-15, M6-17 | Single report captures commands run, evidence, remaining known limitations, and explicit MVP readiness decision. | Not Started |

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
