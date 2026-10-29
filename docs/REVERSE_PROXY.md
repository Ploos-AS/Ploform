# Reverse proxy strategy

**Default recommendation: Caddy.** Ploform remains reverse-proxy-agnostic; Nginx, Traefik, HAProxy and other HTTPS-terminating proxies can be used when they meet the same security contract.

## Current state

The repository has a development-only PostgreSQL Compose service and a Go API bound to loopback by default. There is **no production-ready Docker image or supported public deployment** yet. An internal-only Caddy/API/PostgreSQL Compose scaffold is available for development; it deliberately does not publish ports or provision application database roles. `deploy/caddy/Caddyfile.example` is an illustrative template, not a ready-to-deploy configuration.

## Required proxy contract

- Terminate TLS with a trusted certificate and redirect HTTP to HTTPS.
- Allowlist the intended public hostname; reject unknown Host headers.
- Set `PLOFORM_PUBLIC_ORIGIN` to the exact external HTTPS origin (scheme, hostname and optional port).
- Forward requests to a **loopback-only** Go API on the same host. The same-origin guard accepts non-TLS upstream requests from loopback or exact peer IPs configured in `PLOFORM_TRUSTED_PROXY_IPS` (comma-separated `/32` IPv4 or `/128` IPv6 addresses). The API validates this setting at startup. **This is only a direct-peer check, not a complete container-network security solution.** Docker container IPs can change; the deployment must pin or securely discover the proxy IP and isolate the API network from untrusted containers. Do not expose the upstream port publicly.
- Do not expose PostgreSQL, the internal API port, or administrative endpoints directly.
- Set body-size and request-time limits, rate limits and appropriate access logging without recording form contents or credentials.
- Never treat client-supplied `Forwarded` / `X-Forwarded-*` headers as trustworthy unless a specific trusted-proxy implementation strips and validates them.
- Keep authentication, session management, CSRF and tenant authorization in the application; a proxy does not replace them.

## Future packaging

1. Provide a secure Caddy + app + PostgreSQL deployment with explicit trusted proxy configuration.
2. Provide proxy-free deployment for operators who already terminate HTTPS.
3. Publish tested examples for Nginx, Traefik and HAProxy.
4. Add end-to-end proxy tests before declaring any public deployment supported.

**Do not expose Ploform to the public Internet yet.** Login and other security-critical features remain unfinished.

## Current implementation caveat

The trusted-proxy allowlist is applied to the same-origin guard for unsafe requests, not as a global authentication or host-validation mechanism. The application does not currently provide a production-ready Caddy Compose stack. The public Host allowlist remains the reverse proxy's responsibility. Never rely on `X-Forwarded-For` for caller identity. Tests are written but not executed.

## Internal-only Compose scaffold

Run `PLOFORM_DEV_POSTGRES_PASSWORD=<strong-development-password> docker compose -f deploy/compose.proxy.dev.yaml up --build` from the repository root. The scaffold starts PostgreSQL, an API without database credentials, and an internal Caddy HTTP proxy. It does **not** initialize database migrations or make the API ready for real submissions. It publishes no host ports. The proxy has a fixed Docker IP, which is configured as the only trusted non-loopback direct peer for the API's unsafe-method origin guard. The private network subnet must not overlap existing Docker networks. The internal Caddyfile is deliberately HTTP-only for smoke testing; it must never be reused for public traffic. No deployment test has been run yet.
