# Session foundation

Session token generation and SHA-256 digesting are implemented in
`backend/internal/session`. The database schema stores **only digests**.

This is not a login system. No password verification, session issuance
endpoint, cookie handling, revocation endpoint, or trusted membership
resolver is wired into HTTP yet.

## Requirements for the next implementation

- Authenticate identity before issuing a token.
- Send the raw token only through a Secure, HttpOnly, SameSite cookie.
- Persist only its digest; set short expiration and revoke on logout.
- Enforce CSRF protections for state-changing cookie-authenticated requests.
- Resolve organization membership server-side for each request.
- Never allow a caller-provided organization ID or role to become trusted.
- Rate-limit login and session validation paths.
- Use a dedicated least-privilege authentication database credential;
  the tenant application role must not receive unrestricted session access.

Session helpers are unit-tested, but no full authentication flow has been
implemented or verified.
