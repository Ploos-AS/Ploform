# Ploform

**Open-source forms, surveys and feedback — without artificial limits.**

Ploform is an MIT-licensed, self-hostable forms platform developed by **Ploos AS**. It aims to support independent organizations, brands, workspaces and sites from one deployment, while keeping data and permissions properly isolated.

> **Status: early development, not production-ready.** The Go API now contains database-backed public form routes, submission handling, and authenticated administration routes. These are development features, not a production security guarantee.

## Principles

- Self-hosting first; no mandatory third-party cloud services.
- No license-imposed limits on users, organizations, forms or responses.
- Organization → workspace/project → form hierarchy.
- Explicit tenant boundaries and authorization on every operation.
- Portable configuration and per-organization export.
- Optional webhooks (including n8n) and optional anti-spam providers.
- Privacy by design, with retention and deletion controls.

## Planned stack

Go API · PostgreSQL · SvelteKit admin UI · Docker Compose · OpenAPI.

See [architecture](docs/ARCHITECTURE.md), [roadmap](docs/ROADMAP.md) and [security policy](SECURITY.md).

## Quick start (development scaffold)

```bash
cd backend
go test ./...
go run ./cmd/ploform-api
curl http://localhost:8080/healthz
```

Without a configured database, the API serves health and the contact template; database-backed routes are not enabled. With databases configured, additional development routes are available. **Do not expose this scaffold as a production form service.**

For Ousdal IT, the target is three independent forms (general contact, support, courses) with server-side email notification. **Email delivery is not yet verified as implemented.** Do not embed public forms in production until mail delivery, privacy, and abuse protection have been tested.

## License

MIT — Copyright (c) 2026 Ploos AS. See [LICENSE](LICENSE).

Ploform is independent of any deployment hostname. `forms.ousdal.net` is only the maintainers' internal instance, not the product's public domain.
