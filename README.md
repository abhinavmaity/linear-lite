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

## Current Status

The repository is currently in the planning and architecture-definition phase. The backend technical architecture has been expanded into a detailed implementation reference, and the frontend/user-flow planning has already been outlined. The next stage is execution: turning the documented product, screen flows, schema rules, and API contracts into the actual application.
