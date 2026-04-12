# Milestone 6 Validation Matrix

Date: April 12, 2026  
Scope: MVP journey verification baseline for Milestone 6

## Goal

Map each MVP journey to concrete checks, ownership, and pass criteria. This matrix is the execution reference for M6-02 and feeds M6-03 through M6-18.

## Journey Matrix

| Journey | Coverage Type | Checks | Owner Area | Pass Criteria |
| --- | --- | --- | --- | --- |
| J1: First-time onboarding | Automated + Manual | Register, login, token/session bootstrap, protected route access | Frontend + Backend + QA | User can register/login; protected routes are accessible with token and rejected without token. |
| J2: Daily issue workflow | Automated + Manual | Create issue, list visibility, board visibility, status update, issue detail update, archive/restore | Frontend + Backend + QA | Issue lifecycle completes without contract errors and stays consistent across list/board/detail. |
| J3: Sprint planning flow | Automated + Manual | Create sprint, assign issues, sprint filtering, progress visibility | Frontend + Backend + QA | Sprint-linked issue behavior matches API validation and page states render correctly. |
| J4: Supporting resources CRUD | Automated + Manual | Projects, sprints, labels CRUD plus conflict/dependency handling | Frontend + Backend + QA | Expected success and failure states map to architecture-defined error envelopes and UI handling. |
| J5: Dashboard consistency | Automated + Manual | Dashboard stats load and reflect issue mutations | Backend + QA | Dashboard values update after mutations and exclude archived issues per architecture rules. |
| J6: Team read-only behavior | Manual | Team list/detail rendering, filtering, sorting, no edit affordances | Frontend + QA | Team page remains read-only in MVP and does not expose unsupported write paths. |
| J7: Deployment readiness | Automated + Manual | Frontend build, backend build, backend smoke scripts, full-stack startup path | Platform + Frontend + Backend | Commands are reproducible from documented paths with no missing runtime assumptions. |
| J8: Error and edge handling | Automated + Manual | Validation errors, conflicts, unauthorized access, not-found behavior | Frontend + Backend + QA | Error envelopes and UI states are consistent with contract-defined status codes and messages. |

## Baseline Command Set (M6-03)

From `frontend/`:
- `npm run build`

From `backend/`:
- `GOCACHE=/tmp/go-build go build ./...`

From repo root:
- `./scripts/smoke_issue_workflow.sh`
- `./scripts/smoke_cache.sh`

## Notes

- This matrix intentionally avoids feature expansion and remains aligned to MVP scope.
- Critical browser E2E additions in later Milestone 6 tasks should reference this matrix rather than redefining journey coverage.
