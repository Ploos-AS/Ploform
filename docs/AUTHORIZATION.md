# Authorization design (not yet implemented)

## Identity

Users are global identities. Memberships connect a user to an organization
with a role. Do not infer authorization from an organization ID supplied in
an HTTP path, header or form field.

Planned roles: `owner`, `admin`, `editor`, `viewer`; least privilege by
default. Public form respondents are not organization members.

## Request sequence

1. Authenticate session or API credential on the server.
2. Resolve active membership for the requested organization.
3. Check permission for the requested operation and resource.
4. Begin database transaction with a least-privilege database role.
5. Set transaction-local `app.organization_id` to the **authorized** ID.
6. Execute parameterized SQL; commit or roll back.

No unauthenticated administration endpoints may be added.

## Future schema

- users (identities and secure password / external identity references)
- organization_memberships (organization_id, user_id, role)
- sessions and revocation
- scoped API tokens (hashed, expiring, auditable)
- audit log for administration actions

Security requirements include CSRF protections for cookie sessions, secure
password hashing if local credentials are supported, MFA support, rate
limits, short-lived sessions, secret rotation and authorization integration
tests. Never put access tokens in query parameters.

## Provisioning

Creating organizations is a privileged operation separate from tenant data
access. Do not grant application tenant role unrestricted write privileges
to the organizations table.
