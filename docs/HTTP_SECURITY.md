# HTTP security boundaries

The HTTP authorization middleware verifies an opaque session, then resolves
organization membership and checks a permission. It currently rejects all
unsafe methods and is not mounted on a public API route.

`RequireSameOrigin` is an additional guard for future cookie-authenticated
write endpoints. It requires HTTPS and a matching `Origin` for unsafe
methods. It is **not sufficient alone** for a production deployment:

- Configure a strict external host allowlist; never trust arbitrary Host
  headers or unvalidated forwarded headers.
- For TLS termination at a reverse proxy, set `PLOFORM_PUBLIC_ORIGIN` to
  the exact HTTPS origin (for example `https://forms.example.test`).
  Non-TLS upstream connections are accepted only from loopback clients.
  Bind the Go server to loopback and configure the proxy to reject
  unrecognized public hosts; never expose the upstream port publicly.
- Add synchronizer or double-submit CSRF tokens for state-changing routes.
- Authenticate identity before issuing a session and apply rate limits.
- Set Secure, HttpOnly, host-only cookies and revoke sessions on logout.
- Test login, authorization and CSRF end to end before enabling writes.

The server currently exposes health checks and guarded API endpoints in development mode. No login endpoint exists. Do not expose it publicly before a full deployment security review.
