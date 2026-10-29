# Local development

Requirements: Go 1.23+, Docker Compose, and a PostgreSQL migration runner such as golang-migrate.

1. Copy `.env.example` to `.env` and set a unique password.
2. Run `docker compose up -d postgres`. The database is not published to the host.
3. Run `cd backend && go test ./... && go run ./cmd/ploform-api`.
4. In another terminal, run `curl http://127.0.0.1:8080/healthz`.

The database migration SQL is in `backend/migrations`. It is not automatically applied yet. For development, run the migration from inside the database network or use an ephemeral migration container. **Do not run the down migration on data you need to retain.**

## Important database roles

The initial database superuser is only for migrations. Before enabling tenant-aware application operations, create a separate least-privilege application role (not a table owner, not superuser, no BYPASSRLS), grant only necessary table privileges, and set `app.organization_id` with `SET LOCAL` inside each transaction after authorization. RLS is defense in depth; the API must still check membership and permissions.

## Current limitations

No database connection in the API, authentication, submissions, or UI yet. The compose file intentionally starts **only PostgreSQL**. Never expose the development API publicly.
