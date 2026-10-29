# Development HTTP API

Ploform is not production-ready. No interactive login exists yet.

## Endpoints

- `GET /api/csrf`: requires a valid session and returns `{"csrf_token":"..."}` with `Cache-Control: no-store`.\n- `POST /api/logout`: requires `X-CSRF-Token` obtained from `/api/csrf`, in addition to same-origin verification; revokes the current session and clears its cookie. Requires direct HTTPS with matching `Origin`; current origin middleware does not support TLS-terminating reverse proxies. Do not expose publicly until proxy host validation and CSRF protection are completed.\n- `GET /api/me`: returns `{"user_id":"uuid"}` for an existing valid session, otherwise HTTP 401. It never creates sessions.\n- `GET /healthz`: process liveness.
- `GET /readyz`: checks the tenant database connection only.
- `GET /api/organizations/{organizationID}/workspaces`: requires an authenticated session and verified organization membership.

## Workspace response

```json
{"workspaces":[{"id":"uuid","organization_id":"uuid","slug":"course-registration","name":"Course registration"}]}
```

An empty collection is `{"workspaces":[]}`.

The workspace route is not mounted unless both tenant and authentication database connections are configured. It responds with 401 for missing authentication, 403 for denied access, and 503 for database errors. Organization IDs in paths are requested scopes and must be verified through server-side membership checks.

No write endpoints or login endpoint are enabled yet.

## Forms API v0.1 (incremental)

`GET /api/organizations/{organizationID}/workspaces/{workspaceID}/forms` lists form metadata within a verified organization and requested workspace. Requires a valid session and organization membership. Returns `{ "forms": [] }` for an empty collection. Form fields include `id`, `organization_id`, `workspace_id`, `slug`, `title`, and `status`. Responses are `Cache-Control: no-store`.

This endpoint is **read-only**. Creating drafts, editing schemas, publishing, public rendering and submissions are not implemented or enabled. The endpoint has not yet been integration-tested against PostgreSQL.
