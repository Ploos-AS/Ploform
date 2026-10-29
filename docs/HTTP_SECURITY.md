# HTTP security boundaries

The HTTP authorization middleware verifies an opaque session, then resolves
organization membership and checks a permission. It currently rejects all
unsafe methods and is not mounted on a public API route.

`RequireSameOrigin` is an additional guard for future cookie-authenticated
write endpoints. It requires HTTPS and a matching `Origin` for unsafe
methods. It is **not sufficient alone** for a production deployment:

- Configure a strict external host allowlist; never trust arbitrary Host
  headers or unvalidated forwarded headers.
- Terminate HTTPS at a trusted reverse proxy and pass a verified external
  origin to the application; the current middleware checks `r.TLS` and
  therefore cannot be used unchanged behind TLS termination.
- Add synchronizer or double-submit CSRF tokens for state-changing routes.
- Authenticate identity before issuing a session and apply rate limits.
- Set Secure, HttpOnly, host-only cookies and revoke sessions on logout.
- Test login, authorization and CSRF end to end before enabling writes.

The only public API route currently wired to the server is `GET /healthz`.
