# Integration Roadmap

## Purpose
This roadmap defines the milestone-based path from the current frontend-heavy prototype state to an MVP-complete, integrated full-stack implementation. Use it as the sequencing guide for execution planning, task assignment, and progress tracking. The milestones are sequential, but each milestone also calls out work that can safely proceed in parallel.

## Current Baseline
Linear-lite has clear product scope, frontend planning, design direction, and a detailed backend architecture contract. The frontend includes the core authenticated shell, auth screens, dashboard, issues list, board, issue detail, create issue modal, and integrated supporting resource pages. Frontend development is backend-backed for MVP workflows, and Milestone 6 hardening/readiness tasks are complete.

## Execution Status Update

Completed:
- Milestone 1: Backend Foundation and Runtime Skeleton
- Milestone 2: Database Schema and Core Auth
- Milestone 3: Core Issue Workflow Backend
- Milestone 4: Dashboard and Supporting Resource APIs
- Milestone 5: Frontend Integration Parity Pass
- Milestone 6: Quality, Deployment, and MVP Readiness

Current focus:
- Focused MVP QA and release progression within documented scope

Milestone 5 completion notes (April 10, 2026):
- Removed active mock fallback routing from frontend runtime API calls.
- Completed route parity for dashboard, issues list, board, issue detail, and create issue modal against real backend contracts.
- Completed projects/sprints/labels frontend CRUD parity with backend conflict and validation handling.
- Preserved team page as read-only while improving filter/sort and state handling.
- Added skeleton-loading integration (`boneyard-js`) for major loading routes.
- Added milestone validation artifacts in `docs/Milestone_5_Validation_Report.md` and `docs/Milestone_5_Parity_Completion_Report.md`.

## Milestones

### Milestone 1: Backend Foundation and Runtime Skeleton
**Goal**
Establish the backend runtime, project structure, and shared infrastructure required for all subsequent API implementation.

**Why this milestone exists**
No backend runtime exists in the repository today. This milestone creates the implementation baseline that all domain work will depend on.

**Exact tasks**
- Create the `backend/` structure to match the architecture document, including `cmd/api` and the `internal/` layer directories.
- Add application bootstrap, configuration loading, HTTP server startup, and route registration shell.
- Set up PostgreSQL connection management, migration runner strategy, and local environment conventions.
- Add middleware scaffolding for request IDs, structured logging, panic recovery, CORS, and JWT auth.
- Define shared error model, validation approach, and success/error response envelope conventions.
- Add Docker and backend local-development baseline files and startup expectations.

**Parallelizable tasks**
- Configuration loading, middleware scaffolding, and migration scaffolding can proceed in parallel once the backend directory layout is in place.
- Docker baseline setup can proceed in parallel with server bootstrap if environment variables and port assumptions are aligned.

**Dependencies**
- Backend structure must align with `docs/Technical_Architecture.md` before any domain logic begins.
- Shared response, error, and validation conventions must be fixed before handlers are implemented.

**Definition of done**
- Backend project structure exists and matches the architecture layout.
- The API process boots locally with config loading and route registration shell in place.
- PostgreSQL connectivity and migration strategy are defined and executable.
- Shared middleware scaffolding is present and wired in the canonical request path.
- Docker/local-development baseline exists for backend startup.

**Risks and watchouts**
- Do not let runtime scaffolding drift from the architecture’s layered monolith responsibilities.
- Do not introduce schema management through GORM auto-migration.
- Treat Redis as required infrastructure from the start and align the runtime around it.

### Milestone 2: Database Schema and Core Auth
**Goal**
Implement the database foundation and the authentication flow required to move the frontend auth experience off mock mode.

**Why this milestone exists**
Auth is the smallest complete vertical slice of the backend and unlocks real session handling for the frontend.

**Exact tasks**
- Implement canonical SQL migrations for required extensions, shared functions, triggers, and base schema setup.
- Implement the `users` table, constraints, indexes, and trigger behavior exactly as defined in the architecture document.
- Implement auth repository, service, handler, and route wiring.
- Deliver `POST /auth/register`, `POST /auth/login`, and `GET /auth/me`.
- Implement bcrypt password hashing and JWT access-token issuance/validation.
- Validate frontend login, register, and session-restore flows against the real API.

**Parallelizable tasks**
- User schema migrations and auth service/handler wiring can proceed in parallel once field-level validation and JWT config are fixed.
- Frontend auth integration work can begin once response envelopes and token storage behavior are stable.

**Dependencies**
- Milestone 1 runtime and middleware scaffolding must be complete.
- JWT middleware depends on finalized auth configuration and token parsing rules.

**Definition of done**
- User schema is implemented via SQL migrations and applied successfully.
- Auth endpoints behave according to contract and return the correct envelopes/status codes.
- Frontend login, register, and session bootstrap work against the backend without mock mode for auth.

**Risks and watchouts**
- Enforce architecture limits: access token only, no refresh tokens, no cookie auth.
- Validate email uniqueness case-insensitively.
- Keep password limits aligned with bcrypt constraints and documented validation rules.

### Milestone 3: Core Issue Workflow Backend
**Goal**
Implement the backend slice required for the primary issue-management workflow that powers the existing frontend core screens.

**Why this milestone exists**
The current frontend value is concentrated in issue creation, listing, board movement, detail editing, and archiving. This milestone makes that workflow real end-to-end.

**Exact tasks**
- Implement `projects`, `sprints`, `labels`, `issues`, `issue_labels`, and `issue_activities` schema support required by issue workflows.
- Build repositories, services, handlers, and route wiring for issue list, create, detail, partial update, archive, and restore behavior.
- Deliver selector-support endpoints used by frontend issue flows: `GET /users`, `GET /projects`, `GET /sprints`, and `GET /labels`.
- Deliver issue endpoints required by the frontend: `GET /issues`, `POST /issues`, `GET /issues/:id`, `PUT /issues/:id`, and `DELETE /issues/:id`.
- Implement issue identifier generation, project row locking, archive behavior, activity logging, filtering, sorting, pagination, and archived issue retrieval rules.
- Enforce project/sprint consistency rules during issue creation and updates.

**Parallelizable tasks**
- Selector endpoints can be built in parallel with issue handler/service work once their repository contracts are defined.
- Activity logging and issue-label link management can proceed in parallel with issue list/search implementation.

**Dependencies**
- Milestones 1 and 2 must be complete.
- Project, sprint, and label schemas must exist before issue validation can be enforced correctly.

**Definition of done**
- Create issue, list issues, board status updates, issue detail edits, archive, and restore all work end-to-end against the real backend.
- Frontend core issue workflows no longer require mock data.
- Activity rows are written according to the architecture’s change-tracking rules.

**Risks and watchouts**
- Use the same source of truth for both list and board views.
- Respect archive semantics: `DELETE` archives, `PUT` with `archived: false` restores, `PUT` with `archived: true` is rejected.
- Do not allow sprint/project mismatches on issue create or update.

### Milestone 4: Dashboard and Supporting Resource APIs
**Goal**
Complete the remaining MVP backend surface needed by the dashboard and supporting management screens.

**Why this milestone exists**
After the issue workflow is integrated, the remaining MVP gaps are the dashboard metrics and full CRUD behavior for projects, sprints, and labels.

**Exact tasks**
- Implement `GET /dashboard/stats` according to the architecture-defined metric set.
- Implement remaining project endpoints: `POST /projects`, `GET /projects/:id`, `PUT /projects/:id`, `DELETE /projects/:id`.
- Implement remaining sprint endpoints: `POST /sprints`, `GET /sprints/:id`, `PUT /sprints/:id`, `DELETE /sprints/:id`.
- Implement remaining label endpoints: `POST /labels`, `GET /labels/:id`, `PUT /labels/:id`, `DELETE /labels/:id`.
- Keep the team/users flow read-only from the frontend perspective while still supporting `/users` and `/users/:id`.
- Enforce sprint project-scope rules, one-active-sprint-per-project rules, project deletion rules, and label validation/deletion rules.
- Add Redis-backed cache boundaries only where the architecture allows them.

**Parallelizable tasks**
- Dashboard service work and supporting resource CRUD can proceed in parallel once repositories and shared validation are stable.
- Project, sprint, and label endpoint implementation can be split across separate contributors if repository ownership boundaries stay clear.

**Dependencies**
- Core schemas and issue relations from Milestone 3 must already exist.
- Cache behavior must remain secondary to correct database behavior.

**Definition of done**
- Dashboard metrics work against the real backend.
- Projects, sprints, and labels support their MVP contract-defined CRUD behavior.
- Read-only team frontend flows are backed by stable user endpoints.

**Risks and watchouts**
- Dashboard counts must exclude archived issues.
- Keep permissions simple: all authenticated users can perform MVP actions.
- Do not implement unsupported analytics or admin-only behavior.

### Milestone 5: Frontend Integration Parity Pass
**Goal**
Move the existing frontend from prototype/mock parity to real API parity for the documented MVP user journeys.

**Why this milestone exists**
The frontend already has the right screen footprint, but it still contains mock-backed assumptions and parity gaps that must be resolved once the backend is live.

**Exact tasks**
- Replace mock-backed frontend flows with real API usage milestone by milestone.
- Remove mock-only assumptions in auth, dashboard, issues list, board, issue detail, and create issue modal.
- Close parity gaps in the current frontend: board filters, list sorting UX, richer issue metadata rendering, archive confirmation, and more complete loading/error/empty states.
- Confirm unsupported capabilities remain out of scope: comments, notifications, advanced analytics, multi-workspace, attachments, realtime features.
- Update any UI copy that still frames pages as scaffold/prototype placeholders.
- Validate the user journeys from the frontend planning document against the real backend.

**Parallelizable tasks**
- API integration and UI parity cleanup can proceed in parallel once endpoint contracts are stable.
- Individual page parity passes can be split by route if shared query/state abstractions remain coordinated.

**Dependencies**
- Milestones 2 through 4 must expose stable backend behavior for the routes the frontend consumes.
- Auth/session handling must be stable before protected screen validation begins.

**Definition of done**
- Primary user journeys from `docs/Frontend_Planning.md` run against the real backend without mock data.
- Frontend behavior matches the backend contract for the current MVP screens.
- Prototype/scaffold copy is removed or replaced where inappropriate.

**Risks and watchouts**
- Do not widen scope while closing parity gaps.
- Keep comments and discussion UIs read-only or absent, as defined by the validation report and architecture.
- Validate optimistic board behavior carefully against real failure modes.

### Milestone 6: Quality, Deployment, and MVP Readiness
**Goal**
Harden the integrated product so it is locally runnable, verifiable, and ready for focused QA toward MVP completion.

**Why this milestone exists**
Integration alone is not enough. The system still needs validation coverage, deployment readiness, and a consistency pass across docs and implementation.

**Exact tasks**
- Add backend and frontend validation coverage for core paths.
- Validate architecture rules, edge cases, and failure handling across auth, issues, resources, and dashboard.
- Add Docker Compose or equivalent local full-stack baseline for frontend, backend, and database.
- Perform final documentation consistency pass across README, roadmap, architecture references, and frontend planning references.
- Capture known post-MVP deferrals explicitly so they are not reintroduced during hardening.

**Parallelizable tasks**
- Test coverage work and Docker/dev-environment work can proceed in parallel once core integration is stable.
- Documentation consistency and QA checklist work can proceed in parallel with bug-fix hardening.

**Dependencies**
- Milestones 1 through 5 must be substantially complete.
- Docker readiness depends on settled environment variables, migration execution path, and service ports.

**Definition of done**
- Core flows are validated across frontend and backend.
- Local full-stack startup is documented and reproducible.
- Documentation reflects the integrated state of the repo.
- The project is implementation-grounded, integration-complete for MVP scope, and ready for focused QA.

**Risks and watchouts**
- Avoid treating QA hardening as a feature-expansion phase.
- Keep post-MVP items explicitly deferred.
- Ensure deployment readiness does not introduce undocumented runtime assumptions.

## Dependency Notes
- Milestone 1 is the required foundation for all later work.
- Milestone 2 must complete before protected frontend flows can be validated against the backend.
- Milestone 3 is the unlock for the existing frontend’s highest-value workflows.
- Milestone 4 should not begin until schema ownership, validation, and cache boundaries are already stable.
- Milestone 5 depends on stable route contracts from Milestones 2 through 4.
- Milestone 6 should focus on hardening, not revisiting core product scope.

## Definition of Done Across the Roadmap
The roadmap is complete when all of the following are true:
- The backend exists in code and matches the architecture’s layered monolith responsibilities.
- All required MVP endpoints are implemented and behave according to the documented contract.
- Frontend core flows run against the real backend without mock mode.
- Projects, sprints, labels, users/team, issues, and dashboard support the documented MVP scope.
- Archive behavior, activity logging, filtering, and issue organization rules are enforced.
- Local full-stack development and validation are documented and reproducible.
- README and supporting docs reflect the real implementation state rather than the original planning-only phase.
