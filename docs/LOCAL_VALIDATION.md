## Local validation status

The repository's development Compose smoke test now checks Compose configuration, waits for services, verifies API/Caddy health and database readiness, and requires eight migration-history records. The Dockerfile no longer runs `go mod tidy` during build.

**Not yet validated:** Docker is unavailable in the assistant environment; Go module downloads are blocked by DNS/network restrictions; PostgreSQL and psql are not installed. No end-to-end tests were run. A committed `go.sum` and repeatable offline Go dependency validation are still required. Migration failure rollback, restart idempotency, and cross-tenant isolation must be verified on an operator-controlled Docker host before any public deployment.

For a local test, create two distinct strong secrets under `deploy/secrets/`, set `PLOFORM_DEV_POSTGRES_PASSWORD`, and run `sh deploy/scripts/smoke-test.sh`. Never commit the secret files. The Compose stack does not publish ports.
