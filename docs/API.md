# Development HTTP API

Ploform is not production-ready. No interactive login exists yet.

## Endpoints

- `GET /healthz`: process liveness.
- `GET /readyz`: checks the tenant database connection only.
- `GET /api/organizations/{organizationID}/workspaces`: requires an authenticated session and verified organization membership.

## Workspace response

```json
{"workspaces":[{"id":"uuid","organization_id":"uuid","slug":"course-registration","name":"Course registration"}]}
```

An empty collection is `{"workspaces":[]}`.

The workspace route is not mounted unless both tenant and authentication database connections are configured. It responds with 401 for missing authentication, 403 for denied access, and 503 for database errors. Organization IDs in paths are requested scopes and must be verified through server-side membership checks.

No write endpoints or login endpoint are enabled yet.
