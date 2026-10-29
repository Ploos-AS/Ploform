# Reverse proxy strategy

**Default recommendation: Caddy.** Ploform remains reverse-proxy-agnostic; Nginx, Traefik, HAProxy and other HTTPS-terminating proxies can be used when they meet the same security contract.

## Current state

The repository has a development-only PostgreSQL Compose service and a Go API bound to loopback by default. There is **no production-ready Docker image, Compose application service, Caddy integration or supported public deployment** yet. `deploy/caddy/Caddyfile.example` is an illustrative template, not a ready-to-deploy configuration.

## Required proxy contract

- Terminate TLS with a trusted certificate and redirect HTTP to HTTPS.
- Allowlist the intended public hostname; reject unknown Host headers.
- Set `PLOFORM_PUBLIC_ORIGIN` to the exact external HTTPS origin (scheme, hostname and optional port).
- Forward requests to a **loopback-only** Go API on the same host. Current same-origin guard rejects non-TLS upstream requests from non-loopback remote IPs. A container-to-container proxy network **does not work with this restriction** without a reviewed trusted-proxy mechanism.
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
