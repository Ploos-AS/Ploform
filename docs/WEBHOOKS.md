# Outbound webhooks (design and current status)

Ploform will deliver optional signed HTTPS webhooks to destinations such as n8n, Make, Zapier and custom services. n8n is a receiver, not a Ploform runtime dependency.

## Implemented

`backend/internal/webhooks/signature.go` implements HMAC-SHA256 signatures over the exact transmitted body bytes prefixed by `<unix_timestamp>.`. The sender will use `X-Ploform-Timestamp` and `X-Ploform-Signature: v1=<lowercase hex>`. Receivers should check the timestamp against a five-minute clock-skew window, verify the MAC in constant time, and deduplicate immutable event IDs. Unit tests are written but not executed.

## Required before outbound delivery

- Per-organization, per-destination secrets in an external secrets manager, rotation and revocation.
- Authorization to configure destinations, HTTPS-only allowlisted endpoints and SSRF defenses (DNS rebinding, redirects, private/loopback/link-local addresses).
- Transactional outbox, asynchronous worker, bounded retries with backoff, dead-letter handling and delivery logs without sensitive payloads.
- Stable event IDs, schema versioning, delivery attempt IDs, idempotency and rate limits.
- Safe event payloads: sensitive/confidential answers must not be included by default. Notifications may include a submission ID and authenticated dashboard link only.
- Explicit tenant scoping for destination lookup and event dispatch; no cross-tenant fan-out.

**Not implemented:** webhook destination API, delivery worker, outbox persistence, n8n templates or live webhook transport. No form integration is exposed publicly.
