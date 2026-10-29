# Local checks while GitHub Actions is paused

The GitHub Actions workflow is **manual-only** during the usage limit.
Do not re-enable push or pull_request triggers until the owner approves.

## Go

Requires Go 1.23+ and network access for the first dependency download.

```sh
go version
cd backend
go mod tidy
cd ..
bash scripts/check-local.sh
```

The script checks formatting, runs unit tests and runs go vet. It does
not start a server, install a database or alter dependencies. The
repository does not yet contain a committed go.sum; resolve and commit
that file when dependency fetching is available.

## PostgreSQL

Use a disposable PostgreSQL 17 database. Run as a migration owner,
not as `ploform_app` or `ploform_auth`.

```sh
export PGDATABASE=ploform_test
psql -v ON_ERROR_STOP=1 -f backend/migrations/000001_initial.up.sql
psql -v ON_ERROR_STOP=1 -f backend/migrations/000002_app_role.up.sql
psql -v ON_ERROR_STOP=1 -f backend/migrations/000003_identity.up.sql
psql -v ON_ERROR_STOP=1 -f backend/migrations/000004_sessions.up.sql
psql -v ON_ERROR_STOP=1 -f backend/migrations/000005_auth_role.up.sql
psql -v ON_ERROR_STOP=1 -f backend/migrations/000006_login_challenges.up.sql
psql -v ON_ERROR_STOP=1 -f backend/tests/auth_role_isolation.sql
psql -v ON_ERROR_STOP=1 -f backend/tests/tenant_isolation.sql
psql -v ON_ERROR_STOP=1 -f backend/tests/workspace_access.sql
```

These scripts assume an empty disposable database. They do not test
real login, reverse proxy configuration or a deployed instance.

## CI

To run CI later, use the repository Actions tab and manually start
**Go and PostgreSQL checks**. Do not claim a green build until the
workflow results have been inspected.
