# M0 — Identity, memberships and authorization design

Status: proposal; no production-ready authentication is implied. Tracks #5.

## Core model
- `users`: globally unique user ID, normalized verified email, status, created_at. A user may join multiple organizations.
- `organizations`: ID, name, status, created_at.
- `organization_memberships`: ID, organization_id, user_id, role (`owner|admin|editor|viewer`), status (`active|suspended`), timestamps; UNIQUE (organization_id, user_id).
- `organization_invitations`: ID, organization_id, normalized email, intended role, token_hash, expires_at, accepted_at, revoked_at, invited_by. One-time token only; never store plaintext.
- `audit_events`: ID, organization_id, actor_user_id, action, target_type, target_id, occurred_at. Do not store submission content or secrets.

## Authorization matrix
| Operation | Owner | Admin | Editor | Viewer |
| --- | --- | --- | --- | --- |
| Read forms and responses | yes | yes | yes | yes |
| CSV export | yes | yes | yes | yes |
| Create/edit/publish forms | yes | yes | yes | no |
| Invite and manage non-owner users | yes | yes | no | no |
| Manage owners / organization settings | yes | no | no | no |

All checks happen in Go handlers/services, scoped to the authenticated user and organization. The client-selected organization is never itself proof of authorization. Return not-found for inaccessible tenant resources where appropriate.

## Transactional invariants
- Never remove, suspend or demote the last active owner; serialize owner-changing operations in a transaction and lock the organization row.
- Accept invitations only when unexpired, unrevoked, matching the authenticated verified email, and not already consumed; consume atomically.
- Avoid global uniqueness assumptions for email invitations across organizations.
- Tenant-owned records carry `organization_id`; verify tenant consistency at each relationship boundary.
- Role changes and exports generate audit events.

## Migration strategy
Before writing SQL migrations, inspect existing PostgreSQL schema and migration conventions. Introduce additive migrations with backfill for existing organization data; avoid destructive replacement. Do not introduce a second users/organizations table if they already exist.

## Verification
- Table-driven permission tests for every role and action.
- Integration tests with two organizations, two users, shared member and hostile organization IDs: list, detail, mutate, CSV export and invitation acceptance must not cross boundaries.
- Concurrent demotion/removal of two owners must preserve at least one owner.
- Replayed/expired invitation tokens fail.
- Run tests in CI with PostgreSQL and document results before merge.

## Out of scope
Project-specific ACLs, SSO, analytics/dashboard and production activation. CSV remains the only analysis export; no built-in statistics.
