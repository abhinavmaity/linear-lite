# Backend Task Breakdown

## Source of Truth
This document is the execution-oriented backend implementation checklist for Linear-lite. `docs/Technical_Architecture.md` remains the canonical source of truth for backend behavior, schema, constraints, route contracts, and infrastructure decisions. When this document and the architecture document appear to differ, follow the architecture document and update this checklist to match.

## Implementation Status Snapshot
- Milestones 1 through 4 are complete.
- Milestone 3 (core issue workflow backend) is complete in code, including selector endpoints, issue endpoints, schema support, and archive/restore behavior.
- Milestone 4 (dashboard and supporting resource APIs) is complete in code and validated with cache smoke coverage.
- A reproducible smoke workflow exists at `scripts/smoke_issue_workflow.sh`.
- A reproducible cache workflow exists at `scripts/smoke_cache.sh`.
- CI quality gates are wired via `.github/workflows/ci-validation.yml`.
- Remaining major milestone focus is Milestone 6 quality/deployment hardening and expanded QA coverage.

## Implementation Rules
- Implement the backend as the layered monolith defined by the architecture: router -> middleware -> handlers -> services -> repositories -> PostgreSQL.
- Implement exactly the MVP route surface defined in the architecture. Do not add extra MVP-adjacent endpoints.
- Use SQL migrations as the only schema-management mechanism. Do not use GORM `AutoMigrate`.
- Use JWT bearer authentication with access token only. Do not introduce refresh tokens, cookie auth, or RBAC.
- Treat Redis as required cache infrastructure for the supported runtime while keeping PostgreSQL as the source of truth.
- Enforce business rules in services and structural integrity in the database.
- Return architecture-defined success and error envelopes consistently across all endpoints.
- Preserve single-tenant assumptions and MVP permission simplification: every authenticated user may perform all MVP actions.

## Layer-by-Layer Breakdown

### 1. Bootstrap and Project Structure
Implement the backend repository structure exactly as documented:
- Create `backend/cmd/api/main.go` as the API entrypoint.
- Create `backend/internal/config` for environment/config loading.
- Create `backend/internal/middleware` for request ID, logger, CORS, recovery, and auth middleware.
- Create `backend/internal/handlers` for HTTP request binding and response mapping.
- Create `backend/internal/services` for validation, business rules, transactions, and side effects.
- Create `backend/internal/repositories` for GORM-backed persistence and selective raw SQL.
- Create `backend/internal/models` for GORM table mappings only.
- Create `backend/internal/cache` for Redis-backed cache adapters.
- Create `backend/internal/errors` for shared application/service errors.
- Create `backend/internal/validation` for reusable validation helpers.
- Create `backend/migrations` for SQL-only migrations.

Implement these bootstrap responsibilities:
- Initialize config before any infrastructure objects are created.
- Build the database connection before route registration.
- Register middleware in the canonical order defined by the architecture.
- Register public routes first, then protected routes behind JWT auth.
- Centralize startup logging and fatal-exit behavior in `cmd/api/main.go`.

### 2. Configuration and Runtime Contracts
Implement explicit runtime configuration for:
- Server address and port.
- PostgreSQL DSN or equivalent structured connection fields.
- JWT secret and access-token expiry duration.
- bcrypt hashing cost.
- Redis connection settings.
- Environment name or mode for local/dev/prod behavior.

Enforce these rules:
- Fail fast on missing required config.
- Keep JWT secret mandatory outside test-only setups.
- Require Redis configuration in supported runtime environments and fail fast when it is missing or invalid.
- Keep environment handling simple and explicit rather than magic defaults.
- Document the expected local `.env` or compose-level variable set when implementation begins.

### 3. Database and Migrations
Implement SQL-only migrations in dependency order.

Implement foundational migrations for:
- `pgcrypto` extension.
- Shared `set_updated_at()` trigger function.
- Table-level `updated_at` triggers where required.

Implement canonical tables in this order:
1. `users`
2. `projects`
3. `sprints`
4. `labels`
5. `issues`
6. `issue_labels`
7. `issue_activities`

Implement all architecture-defined:
- Primary keys, foreign keys, and `ON DELETE`/`ON UPDATE` rules.
- Check constraints.
- Functional unique indexes.
- Search support for issues.
- Partial unique indexes such as one active sprint per project.
- Archive-related constraints for issues.

Enforce these migration rules:
- Do not let GORM create or mutate schema at runtime.
- Keep `up` and `down` migrations paired.
- Keep schema authority in SQL, not in model struct tags.
- Validate the migration order against table dependencies before finalizing.

### 4. Domain Models and Repository Layer
Implement GORM model mappings for each canonical table:
- `users`
- `projects`
- `sprints`
- `labels`
- `issues`
- `issue_labels`
- `issue_activities`

Repository ownership must be resource-specific:
- Auth/User repository behavior for user lookup and auth-linked queries.
- Project repository for project CRUD, counts, and detail assembly.
- Sprint repository for sprint CRUD, detail assembly, and project-scoped validations.
- Label repository for label CRUD and usage checks.
- Issue repository for issue list/search/detail/create/update/archive flows.
- Activity repository for issue activity writes and timeline reads.

Use selective raw SQL where the architecture expects it or where PostgreSQL-specific behavior is clearer and safer than ORM composition:
- Issue search and full-text search logic.
- Project row locking and issue identifier generation.
- Any query that requires explicit PostgreSQL features, locking, or clearer SQL semantics than nested ORM composition.

Implement repository behavior for:
- Pagination and sorting.
- Repeated query parameter filters.
- Association loading for detail and collection responses.
- Scoped transactions provided by the service layer.
- Hydrated return shapes that support the documented response envelopes.

### 5. Middleware
Implement these middleware modules:
- Request ID middleware.
- Structured logging middleware.
- Recovery middleware.
- CORS middleware.
- JWT auth middleware.

Enforce these route-protection rules:
- Public routes: `POST /auth/register`, `POST /auth/login`.
- Protected routes: every other MVP endpoint.
- Auth middleware must parse `Authorization: Bearer <token>` only.
- Invalid, missing, or expired tokens must return `401 unauthorized`.

Implement middleware behavior so that:
- Request IDs are attached early and available to logs and error responses where applicable.
- Panic recovery prevents process crashes and returns `500 internal_error`.
- Logging captures method, path, status, duration, and request ID.
- CORS behavior is explicit and environment-compatible.

### 6. Validation and Error Mapping
Implement shared validation helpers for:
- UUID path params.
- Pagination query params.
- Sort fields and sort order.
- Enum values.
- String length and format checks.
- Date parsing.
- Optional/null field handling.
- Repeated query parameter arrays.

Enforce service-level business rule validation for:
- Unique email and project key rules.
- Sprint project ownership rules.
- One-active-sprint-per-project rule.
- Label uniqueness and hex color rules.
- Issue project/sprint alignment.
- Archive and restore semantics.
- Deletion restrictions for dependent resources.

Implement error taxonomy to HTTP mapping:
- `400 validation_error` for invalid input or invalid business-rule requests.
- `401 unauthorized` for auth failures.
- `403 unauthorized` only if a future permissions path is introduced; do not invent MVP-only RBAC checks.
- `404 not_found` for missing resources.
- `409 conflict` for uniqueness and state conflicts defined by the architecture.
- `500 internal_error` for unexpected failures.

Return envelopes consistently:
- Success collection: `{ items, pagination }`
- Success resource: `{ data: ... }`
- Error: `{ error: { code, message, fields?, request_id? } }`

### 7. Service Layer by Domain

#### Auth Service
Implement responsibilities:
- Register users.
- Authenticate users.
- Return the current authenticated user.
- Issue JWT access tokens.

Enforce business rules:
- Email must be valid, trimmed, max 255, and unique case-insensitively.
- Password must respect documented min/max length.
- Name must be present and trimmed.

Use repositories:
- User repository for email lookups and inserts.

Transaction needs:
- Register can run in a simple create-user transaction.

Activity logging requirements:
- None.

Validation rules:
- Reject invalid email/password/name input with `400`.
- Return `409` for duplicate email on register.
- Return `401` for bad credentials on login.

Edge cases:
- Normalize email before lookup and uniqueness checks.
- Do not leak whether email exists on bad login beyond `401 unauthorized`.

#### User Service
Implement responsibilities:
- List users for selectors and team page.
- Return single-user detail with issue statistics.

Enforce business rules:
- Read-only behavior for MVP from the frontend perspective.

Use repositories:
- User repository.
- Issue repository or reporting query support for statistics aggregation.

Transaction needs:
- None for reads.

Activity logging requirements:
- None.

Validation rules:
- Validate pagination, sorting, and UUID path params.

Edge cases:
- Keep search behavior simple: name/email substring search.
- Return `404` for unknown user IDs.

#### Project Service
Implement responsibilities:
- List projects.
- Create, fetch, update, and delete projects.
- Return issue counts and active sprint summary in project responses.

Enforce business rules:
- `key` must match `^[A-Z0-9]{2,10}$`.
- `key` can change only if the project has zero issues.
- Project deletion is blocked when any issue or sprint exists.

Use repositories:
- Project repository.
- Sprint repository for detail hydration if needed.
- Issue repository for counts/dependency checks.

Transaction needs:
- Create and update should be transactional if multiple writes or validations are combined.
- Delete should validate dependencies before deletion in one consistent flow.

Activity logging requirements:
- None.

Validation rules:
- Validate name length, description length, and key format.
- Return `409` for duplicate key or forbidden key change.

Edge cases:
- Do not allow client control of `next_issue_number`.
- Keep response hydration consistent across list and detail routes.

#### Sprint Service
Implement responsibilities:
- List sprints.
- Create, fetch, update, and delete sprints.
- Return issue counts and parent project in sprint detail.

Enforce business rules:
- Sprints are always project-scoped.
- Only one active sprint per project.
- `end_date` must be greater than or equal to `start_date`.
- Delete only if sprint is not active and has no issue references.

Use repositories:
- Sprint repository.
- Project repository.
- Issue repository for counts and dependency checks.

Transaction needs:
- Create/update/delete flows should be transactional when status checks and writes must remain consistent.

Activity logging requirements:
- None.

Validation rules:
- Validate dates, status enum, and project existence.
- Return `409` when active-sprint conflicts occur.

Edge cases:
- Allow overlapping non-active sprints in MVP.
- Reject create/update when referenced project does not exist.

#### Label Service
Implement responsibilities:
- List labels.
- Create, fetch, update, and delete labels.
- Return usage count on label detail.

Enforce business rules:
- Label names must be unique case-insensitively.
- Color must match `#RRGGBB`.
- Delete only when label is unused.

Use repositories:
- Label repository.
- Issue repository or join-table checks for usage validation.

Transaction needs:
- Create/update/delete can use simple transactions where uniqueness/use checks and mutation must stay aligned.

Activity logging requirements:
- None directly for label CRUD.

Validation rules:
- Validate name length, color format, and description length.
- Return `409` for duplicate name or delete-while-used conflicts.

Edge cases:
- Keep usage count queries consistent with archived and non-archived issue associations as defined by repository logic.

#### Issue Service
Implement responsibilities:
- List issues for list and board views.
- Create issues with identifier generation.
- Fetch issue detail with relations and activity timeline.
- Partially update issues.
- Archive issues.
- Restore archived issues.

Enforce business rules:
- Title required, 1-500 chars.
- Description optional, max documented length.
- Status and priority must be valid enums.
- `project_id` must exist.
- `sprint_id`, if present, must exist and belong to the issue’s project.
- `assignee_id`, if present, must exist.
- `label_ids`, if present, must be distinct existing labels.
- `DELETE /issues/:id` archives instead of deleting.
- `PUT /issues/:id` may restore only with `{ "archived": false }`.
- `PUT /issues/:id` with `archived: true` must return `400`.
- Each changed scalar field must write one activity row.
- Label set replacement must write added/removed label activity entries.

Use repositories:
- Issue repository.
- Project repository.
- Sprint repository.
- Label repository.
- User repository.
- Activity repository.

Transaction needs:
- Create must be transactional around project row locking, identifier generation, issue insert, issue-label inserts, activity inserts, and project sequence increment.
- Update must be transactional for field updates, label rewrites, restore operations, and activity writes.
- Archive must be transactional for archive field mutation and activity write.

Activity logging requirements:
- Write `created` activity on create.
- Write per-field updates on scalar changes.
- Write `label_added` and `label_removed` when label membership changes.
- Write `archived` on archive.
- Write `restored` on restore.

Validation rules:
- Validate all filters, repeated query params, sort fields, include-archived flags, and body field shapes.
- Return `404` for missing resources.
- Return `409` for identifier collisions or rule conflicts the architecture marks as conflicts.

Edge cases:
- Board and list must read from the same issue source of truth.
- Archived issues are excluded by default and only included when explicitly requested.
- Repeated archive of an already archived issue returns `204` without further mutation.
- Project changes with mismatched sprint association must be rejected unless sprint is cleared or replaced.

#### Dashboard Service
Implement responsibilities:
- Return dashboard metrics and recent activity.

Enforce business rules:
- Dashboard counts exclude archived issues.
- Recent activity should reflect issue activity records and active issue scope.

Use repositories:
- Issue repository.
- Activity repository.
- Sprint repository if active sprint summary is not already available through issue/project aggregations.

Transaction needs:
- None for reads.

Activity logging requirements:
- None.

Validation rules:
- No request body or query validation beyond auth.

Edge cases:
- Return stable empty states when no active sprint or recent activity exists.
- Keep metric calculations aligned with architecture-defined fields only.

### 8. Handler and Route Delivery
Implement handlers so they only:
- Parse path/query/body input.
- Call the appropriate service.
- Map service results and errors to HTTP responses.
- Never perform raw SQL or business-rule decisions in the handler layer.

#### Auth Routes
| Method | Path | Auth | Handler responsibility | Service | Critical inputs | Success | Major failures |
| --- | --- | --- | --- | --- | --- | --- | --- |
| POST | `/auth/register` | No | Bind registration payload and return token + user envelope | Auth Service | `name`, `email`, `password` | `201` | `400`, `409`, `500` |
| POST | `/auth/login` | No | Bind credentials and return token + user envelope | Auth Service | `email`, `password` | `200` | `400`, `401`, `500` |
| GET | `/auth/me` | Yes | Resolve authenticated user from auth context | Auth Service | auth context only | `200` | `401`, `500` |

#### User Routes
| Method | Path | Auth | Handler responsibility | Service | Critical inputs | Success | Major failures |
| --- | --- | --- | --- | --- | --- | --- | --- |
| GET | `/users` | Yes | Bind list query params and return paginated users | User Service | `page`, `limit`, `search`, `sort_by`, `sort_order` | `200` | `400`, `401`, `500` |
| GET | `/users/:id` | Yes | Bind UUID path param and return user detail + stats | User Service | `id` | `200` | `400`, `401`, `404`, `500` |

#### Project Routes
| Method | Path | Auth | Handler responsibility | Service | Critical inputs | Success | Major failures |
| --- | --- | --- | --- | --- | --- | --- | --- |
| GET | `/projects` | Yes | Bind list query params and return paginated projects | Project Service | `page`, `limit`, `search`, `sort_by`, `sort_order` | `200` | `400`, `401`, `500` |
| POST | `/projects` | Yes | Bind create payload and return project detail | Project Service | `name`, `description`, `key` | `201` | `400`, `401`, `409`, `500` |
| GET | `/projects/:id` | Yes | Bind UUID path param and return project detail | Project Service | `id` | `200` | `400`, `401`, `404`, `500` |
| PUT | `/projects/:id` | Yes | Bind partial update payload and return updated project | Project Service | `id`, optional `name`, `description`, `key` | `200` | `400`, `401`, `404`, `409`, `500` |
| DELETE | `/projects/:id` | Yes | Bind UUID path param and delete eligible project | Project Service | `id` | `204` | `400`, `401`, `404`, `409`, `500` |

#### Sprint Routes
| Method | Path | Auth | Handler responsibility | Service | Critical inputs | Success | Major failures |
| --- | --- | --- | --- | --- | --- | --- | --- |
| GET | `/sprints` | Yes | Bind list query params and return paginated sprints | Sprint Service | `page`, `limit`, `project_id`, `status`, `search`, `sort_by`, `sort_order` | `200` | `400`, `401`, `500` |
| POST | `/sprints` | Yes | Bind create payload and return sprint detail | Sprint Service | `name`, `description`, `project_id`, `start_date`, `end_date`, optional `status` | `201` | `400`, `401`, `404`, `409`, `500` |
| GET | `/sprints/:id` | Yes | Bind UUID path param and return sprint detail | Sprint Service | `id` | `200` | `400`, `401`, `404`, `500` |
| PUT | `/sprints/:id` | Yes | Bind partial update payload and return updated sprint | Sprint Service | `id`, optional `name`, `description`, `start_date`, `end_date`, `status` | `200` | `400`, `401`, `404`, `409`, `500` |
| DELETE | `/sprints/:id` | Yes | Bind UUID path param and delete eligible sprint | Sprint Service | `id` | `204` | `400`, `401`, `404`, `409`, `500` |

#### Label Routes
| Method | Path | Auth | Handler responsibility | Service | Critical inputs | Success | Major failures |
| --- | --- | --- | --- | --- | --- | --- | --- |
| GET | `/labels` | Yes | Bind list query params and return paginated labels | Label Service | `page`, `limit`, `search`, `sort_by`, `sort_order` | `200` | `400`, `401`, `500` |
| POST | `/labels` | Yes | Bind create payload and return label resource | Label Service | `name`, `color`, optional `description` | `201` | `400`, `401`, `409`, `500` |
| GET | `/labels/:id` | Yes | Bind UUID path param and return label detail | Label Service | `id` | `200` | `400`, `401`, `404`, `500` |
| PUT | `/labels/:id` | Yes | Bind partial update payload and return updated label | Label Service | `id`, optional `name`, `color`, `description` | `200` | `400`, `401`, `404`, `409`, `500` |
| DELETE | `/labels/:id` | Yes | Bind UUID path param and delete eligible label | Label Service | `id` | `204` | `400`, `401`, `404`, `409`, `500` |

#### Issue Routes
| Method | Path | Auth | Handler responsibility | Service | Critical inputs | Success | Major failures |
| --- | --- | --- | --- | --- | --- | --- | --- |
| GET | `/issues` | Yes | Bind list query params and return paginated issues | Issue Service | `page`, `limit`, `sort_by`, `sort_order`, `search`, repeated `status`, repeated `priority`, `assignee_id`, `project_id`, `sprint_id`, repeated `label_id`, `label_mode`, `include_archived` | `200` | `400`, `401`, `500` |
| POST | `/issues` | Yes | Bind create payload and return issue detail | Issue Service | `title`, `description`, `status`, `priority`, `project_id`, `sprint_id`, `assignee_id`, `label_ids` | `201` | `400`, `401`, `404`, `409`, `500` |
| GET | `/issues/:id` | Yes | Bind UUID path param and include-archived query flag, then return issue detail | Issue Service | `id`, optional `include_archived` | `200` | `400`, `401`, `404`, `500` |
| PUT | `/issues/:id` | Yes | Bind partial update or restore payload and return issue detail | Issue Service | `id`, optional `title`, `description`, `status`, `priority`, `project_id`, `sprint_id`, `assignee_id`, `label_ids`, `archived` | `200` | `400`, `401`, `404`, `409`, `500` |
| DELETE | `/issues/:id` | Yes | Bind UUID path param and archive issue | Issue Service | `id` | `204` | `400`, `401`, `404`, `500` |

#### Dashboard Route
| Method | Path | Auth | Handler responsibility | Service | Critical inputs | Success | Major failures |
| --- | --- | --- | --- | --- | --- | --- | --- |
| GET | `/dashboard/stats` | Yes | Return dashboard metrics and recent activity | Dashboard Service | auth context only | `200` | `401`, `500` |

### 9. Caching and Redis Infrastructure
Implement cache behavior as Redis-backed wrappers around selected read endpoints only.

Keep Redis required:
- The application startup must fail if Redis is unavailable in supported runtime environments.
- Cache misses must never break request correctness once Redis is connected.
- Cache invalidation must follow successful writes; failed writes must not invalidate as though they succeeded.

Identify candidate read endpoints for later cache wrapping:
- `GET /labels`
- `GET /users`
- `GET /projects`
- `GET /projects/:id`
- `GET /sprints`
- `GET /sprints/:id`
- `GET /dashboard/stats`

Apply invalidation expectations:
- Project writes invalidate relevant project cache and dependent issue/sprint views.
- Sprint writes invalidate sprint, project, dashboard, and affected issue views where applicable.
- Label writes invalidate label and affected issue cache.
- Issue writes invalidate affected project, sprint, dashboard, label, and issue-related cached reads as needed.

### 10. Deployment and Local Environment
Implement the local runtime baseline so frontend, backend, and database can run together.

Deliver these expectations:
- Docker and Docker Compose support for local full-stack startup.
- Explicit relationship between frontend container, backend container, PostgreSQL, and Redis.
- Clear startup expectation for applying migrations during local bootstrapping.
- Environment variables documented for backend startup, frontend API targeting, and local database connectivity.

Enforce these deployment rules:
- Keep PostgreSQL as the primary transactional store.
- Require Redis in the default local path and document it as part of the standard stack.
- Make migration execution part of the documented local setup flow.
- Keep frontend/backend port mapping aligned with the architecture and current frontend expectations.

## Cross-Cutting Rules
- Implement all 26 architecture-defined routes and no extra MVP routes unless the architecture document is first updated.
- Use the same issue source of truth for list and board views.
- Keep archived issues excluded by default unless explicitly requested.
- Keep comments, attachments, notifications, realtime updates, integrations, and multi-workspace support out of scope.
- Keep permissions simplified to authenticated-user access for all MVP actions.
- Keep response shapes stable so the current frontend service layer can integrate without ad hoc shape translation.
- Prefer deterministic, explicit service logic over hidden ORM behavior.

## Acceptance Checklist
- All required backend directories and bootstrap files exist.
- All required config inputs are defined and validated.
- All SQL migrations exist, apply cleanly, and match the architecture-defined schema.
- GORM `AutoMigrate` is not used.
- All 26 MVP endpoints are implemented.
- Auth works end-to-end with the frontend.
- Issues support list and board off the same source of truth.
- Issue archive and restore behavior works exactly as documented.
- Issue activity logging works for create, update, label changes, archive, and restore.
- Project, sprint, and label validation/deletion rules are enforced.
- Dashboard metrics exclude archived issues and return the documented fields only.
- Local full-stack development works with frontend, backend, and PostgreSQL.
