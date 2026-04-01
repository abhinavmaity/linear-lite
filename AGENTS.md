# AGENTS.md

## Purpose

This file defines how AI coding agents should operate in the `linear-lite` repository.

## Project Snapshot

- Product: Linear-lite (lightweight issue/project tracker MVP).
- Current implementation status: frontend exists; backend is primarily architecture/planning.
- Primary stack in code today: React + TypeScript + Vite + Zustand + TanStack Query.

## Source of Truth (Read First)

1. `README.md`
2. `docs/Objective.md`
3. `docs/Frontend_Planning.md`
4. `frontend/design.md` (UI design direction and patterns)
5. `docs/Technical_Architecture.md`
6. `docs/Frontend_Validation_Report.md` (if present and relevant)

When implementation details conflict with assumptions, trust documented architecture and planning docs.

## Repository Layout

- `frontend/`: React app (active implementation area)
- `docs/`: product, UX flow, and architecture references

## Working Rules for Agents

1. Keep scope MVP-focused. Do not add out-of-scope features unless explicitly requested.
2. Prefer small, reviewable changes over large refactors.
3. Preserve existing patterns and naming conventions in `frontend/src`.
4. Do not introduce new dependencies unless they are clearly necessary.
5. Avoid speculative backend implementation not grounded in documented contracts.
6. Never revert user-authored changes you did not create.
7. If uncertain, implement the smallest safe change and document assumptions.

## Frontend Standards

- Language: TypeScript.
- Framework: React functional components.
- State:
  - server state via TanStack Query
  - UI/session/local state via Zustand stores
- API calls should go through `frontend/src/services/*`.
- Keep pages in `frontend/src/pages/*` and reusable UI in `frontend/src/components/*`.
- UI styling and layout decisions should follow `frontend/design.md`.
- Favor clear, predictable UX over visual complexity.

## Commands

Run from `frontend/`:

- Install deps: `npm install`
- Dev server: `npm run dev`
- Build/type-check: `npm run build`
- Preview build: `npm run preview`

## Change Validation

Before finishing a task, agents should:

1. Run `npm run build` in `frontend/` for compile/type validation.
2. Confirm changed screens/components load without obvious runtime errors.
3. Summarize:
   - files changed
   - behavior change
   - any assumptions or follow-ups

If validation cannot be run, explicitly state what was not verified.

## Commit Guidance

- Use focused commits with clear messages.
- Keep unrelated edits out of scope.
- If a task is exploratory, prefer draft PR-style notes in the handoff summary.

## Non-Goals (Unless Explicitly Requested)

- Full backend implementation beyond existing architecture docs
- New infrastructure/tooling migrations
- Broad UI redesign not tied to a concrete requirement
- Premature optimization
