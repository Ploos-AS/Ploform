# Database access

The `Store` abstraction scopes workspace queries to a transaction with
`SET LOCAL app.organization_id` (implemented via `set_config(..., true)`).

**Security invariant:** `orgID` must be resolved from an authenticated,
authorized server-side membership. Passing an organization ID from a URL or
header without authorization is a cross-tenant vulnerability.

The database connection must use a least-privilege `ploform_app` role.
The schema migration role must be different and must not serve HTTP traffic.

The role migration deliberately does **not** set a password. Configure a
secret out of band before connecting the application. Organization creation
and membership provisioning are deliberately not implemented yet.

Future integration tests must verify that tenant A cannot select, update,
insert or delete tenant B's rows, even with a forged organization_id.
