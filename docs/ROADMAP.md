# Roadmap

This roadmap describes intent, not shipped functionality.

## Phase 0 — foundations
- [x] Project README, MIT license and architecture notes
- [x] Minimal Go health endpoint and test
- [ ] CI checks, linting, dependency review and container build
- [ ] PostgreSQL migrations and local development Compose stack

## Phase 1 — safe MVP
- [ ] Organizations, workspaces, memberships and authorization
- [ ] Forms with versioned schemas and validation
- [ ] Public submission API with anti-abuse controls
- [ ] Tenant-scoped submissions and basic export
- [ ] Signed webhook outbox with retries
- [ ] Retention and deletion controls
- [ ] Security and isolation tests

## Phase 2 — usability
- [ ] SvelteKit admin UI and visual field editor
- [ ] Embedded and hosted forms
- [ ] Survey question types and conditional logic
- [ ] Email notifications through optional integrations
- [ ] Accessibility and localization

## Phase 3 — advanced
- [ ] Anonymous surveys and privacy-preserving analytics
- [ ] File uploads with safe storage policies
- [ ] Per-organization portable backup/export
- [ ] Additional anti-spam providers and integrations

No feature should require a proprietary license tier.
