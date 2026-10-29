# PostgreSQL isolation tests

These SQL checks validate row-level security under the restricted `ploform_app` role.
They do **not** test authentication, membership authorization or public API behavior.

Use a disposable database. Apply both migrations first. Run with psql as
the migration owner (or another role authorized to SET ROLE ploform_app):

```sh
psql "$TEST_DATABASE_URL" -v ON_ERROR_STOP=1 -f backend/tests/tenant_isolation.sql
```

The test seeds two organizations and rolls back at the end. It asserts that
tenant A cannot read, insert, update or delete tenant B's workspace.

**Important:** `SET ROLE` and `set_config` in this test are controlled
by a privileged test runner. The HTTP application must never expose either
operation to an untrusted caller.

This script is not yet wired into CI; run it against a disposable PostgreSQL
database before relying on the isolation policy.
