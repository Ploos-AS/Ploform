# Authentication database privilege regression

Run after migrations 000001–000005 using the database migration owner:

```sh
psql -v ON_ERROR_STOP=1 -f backend/tests/auth_role_isolation.sql
```

The test creates temporary identity and session rows, switches to
`ploform_auth`, confirms it can read the authentication tables, and
checks that it cannot read tenant workspaces or create users. The outer
transaction is rolled back.

This tests PostgreSQL grants, not HTTP authentication or cookie issuance.
