# Architecture (draft)

## Product vs deployment

Ploform is a reusable open-source product owned by Ploos AS. Hostnames such as `forms.ousdal.net` are **deployment-specific** and must never be hard-coded into the product.

## Components

- **Go API:** public form submission and authenticated administration endpoints (planned).
- **PostgreSQL:** persistent tenant-aware data.
- **SvelteKit:** form builder and administration (planned).
- **Worker / outbox:** reliable asynchronous webhook delivery (planned).
- **Docker Compose:** local and self-hosted deployment (planned).

## Data hierarchy

`Organization → Workspace/Project → Form → Version → Submission`

Organizations represent security and data-ownership boundaries; projects group forms belonging to a brand or site. Membership and role assignments are scoped explicitly. Every tenant-owned record must carry an organization identifier; queries must be scoped and checked server-side.

### Isolation requirements

- Never trust tenant identifiers from client input without authorization.
- Use database constraints and integration tests for cross-tenant access.
- Consider PostgreSQL row-level security as defense in depth, not a substitute for application checks.
- Organization-scoped export and deletion; avoid cross-tenant backups as a migration mechanism.
- Secrets and webhook signing keys scoped per organization; redact sensitive logs.
- Document processor/controller responsibilities for independently operated organizations.

## Security requirements before public submissions

Request-size limits, schema validation, abuse prevention, server-side anti-bot verification, origin policy, rate limits, safe file handling (if enabled), retention controls, HTTPS, and auditability. CORS is not an authorization mechanism.

## Integration

Webhooks are optional, signed, queued and retried. n8n is one possible receiver, not a dependency. CAPTCHA/Turnstile is a pluggable optional provider.

## Schema evolution

Version form definitions so historic submissions remain interpretable. Define a stable export format and version it separately from the internal database schema.
