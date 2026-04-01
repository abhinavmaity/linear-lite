# Frontend Validation Report

This report validates the `temp-ui` mock screens against [`Frontend_Planning.md`](/Users/abhinavmaity/code/linear-lite/docs/Frontend_Planning.md) and [`Technical_Architecture.md`](/Users/abhinavmaity/code/linear-lite/docs/Technical_Architecture.md). The technical architecture is the source of truth whenever the documents differ.

## Validation Rules

- Light and dark variants are treated as one logical screen unless they diverge structurally.
- Findings are classified as:
  - `Blocker`: cannot be implemented faithfully against the backend contract
  - `Gap`: missing MVP behavior from planning or architecture
  - `Enhancement`: useful visual/design idea that does not block MVP

## Cross-Screen Findings

- `Blocker`: Issue detail mockups include comment entry and reply behavior, but comments are out of scope in the architecture. The implemented screen must treat this area as read-only activity history.
- `Gap`: Several screens imply richer analytics, notifications, settings, and help interactions that are not defined in the MVP API contract. These should stay decorative or disabled until a later phase.
- `Gap`: The planning doc calls for specific filter sets and loading/error states on issues screens; mockups capture the layout direction but not all required states.
- `Enhancement`: The paired light/dark visual direction is strong and can be preserved via a shared token-based theme system rather than duplicated page implementations.

## Screen Audit

### 1. Login

- Route/purpose: matches `/login` and authentication entry intent
- Required components: present for primary auth flow, but the light mock uses `username` instead of architecture-required `email`
- Architecture compatibility:
  - `Blocker`: login must submit `{ email, password }` to `POST /api/v1/auth/login`
  - `Enhancement`: branding, manifesto copy, and secondary actions can stay if they do not alter the core flow
- Reusable patterns: auth card layout, marketing split layout, branded submit CTA

### 2. Register

- Route/purpose: matches `/register` and account creation intent
- Required components: name, email, password, confirm password are represented
- Architecture compatibility:
  - `Gap`: confirm-password is frontend-only validation, not part of the backend payload
  - `Enhancement`: long-form brand copy and consent messaging are acceptable as presentation-only additions
- Reusable patterns: auth card layout, field grouping, inline validation

### 3. Dashboard

- Route/purpose: matches `/dashboard`
- Required components: overview cards, active sprint cue, recent activity, create issue entry are present
- Architecture compatibility:
  - `Gap`: dashboard data must be constrained to `GET /api/v1/dashboard/stats`
  - `Gap`: unsupported mock metrics and charts must become decorative only or be removed
- Reusable patterns: stat cards, activity list, shell layout

### 4. Issues List View

- Route/purpose: matches `/issues`
- Required components: search, table, pagination, create issue entry, view navigation are present
- Architecture compatibility:
  - `Gap`: implemented UI still needs status, priority, assignee, label, sprint, and sorting controls aligned to `GET /api/v1/issues`
  - `Gap`: loading, empty, and error states are not fully represented in the mock
- Reusable patterns: data table, filter bar, view switch affordances

### 5. Issues Board View

- Route/purpose: matches `/board`
- Required components: status columns, cards, create issue entry, shell navigation are present
- Architecture compatibility:
  - `Gap`: board columns must map to backend statuses from the issues API
  - `Gap`: drag-and-drop requires optimistic `PUT /api/v1/issues/:id` updates with rollback on failure
- Reusable patterns: kanban column, issue card, column counts

### 6. Issue Detail

- Route/purpose: matches `/issues/:id`
- Required components: main content, sidebar metadata, activity history, editable fields are represented
- Architecture compatibility:
  - `Blocker`: comment composer/reply UI must not be implemented as writable behavior
  - `Gap`: autosave and field-level updates need to align to partial `PUT /api/v1/issues/:id`
  - `Gap`: archive action must use `DELETE /api/v1/issues/:id`
- Reusable patterns: metadata sidebar, activity timeline, editable field sections

### 7. Create Issue Modal

- Route/purpose: matches create issue modal trigger from any screen
- Required components: title and secondary field inputs are represented
- Architecture compatibility:
  - `Gap`: create payload must use `project_id`, `sprint_id`, `assignee_id`, and `label_ids` exactly as documented
  - `Gap`: required title validation is implied visually but should be explicit in implementation
- Reusable patterns: modal shell, compact form layout

### 8. Projects Page

- Route/purpose: matches `/projects`
- Required components: project list/cards and create affordance are present
- Architecture compatibility:
  - `Gap`: initial implementation should focus on selector/list support and a scaffolded route before deeper CRUD
  - `Enhancement`: analytics side panels can remain decorative until backed by contract data
- Reusable patterns: resource cards, section metrics, shell layout

### 9. Sprints Page

- Route/purpose: matches `/sprints`
- Required components: sprint list/cards, active sprint emphasis, issue count cues are present
- Architecture compatibility:
  - `Gap`: sprint state must reflect `planned`, `active`, `completed`
  - `Gap`: first pass should prioritize list rendering and selector support before full CRUD flows
- Reusable patterns: status cards, resource list blocks

### 10. Labels Management

- Route/purpose: matches `/labels`
- Required components: labels grid and create/edit affordance are present
- Architecture compatibility:
  - `Gap`: implemented CRUD must respect label name uniqueness and hex color format
  - `Enhancement`: visual label gallery is stronger than the minimum MVP and worth preserving
- Reusable patterns: label pills, grid cards, create tile

### 11. Team Page

- Route/purpose: matches `/team`
- Required components: member listing and profile summary are present
- Architecture compatibility:
  - `Gap`: team page must stay read-only in MVP
  - `Enhancement`: rich member presentation is acceptable if it stays purely informational
- Reusable patterns: avatar cards, member stat rows

## Implementation Guidance

- Build one shared React component system with runtime theme switching instead of duplicating light/dark screens.
- Prioritize the core user flow first: auth, dashboard, issues list, board, issue detail, and create issue.
- Preserve the Bauhaus-inspired visual language, but keep all interactive behavior bound strictly to the architecture-defined contracts.
