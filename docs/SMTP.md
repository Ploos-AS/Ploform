# SMTP notifications (development)

The `backend/internal/notify` package provides an opt-in SMTP transport with mandatory STARTTLS, TLS certificate verification, SMTP authentication, server-side sender/recipient addresses and header-injection checks.

**Not yet wired into the submission workflow.** No emails are sent by Ploform until a server-side integration and reliable notification queue are implemented and tested.

## Planned integration
1. Store the submission transactionally, then enqueue a notification job with a stable ID.
2. A background worker sends mail with retry/backoff and records delivery outcomes. Do not fail or duplicate the form submission because SMTP is temporarily unavailable.
3. Configure SMTP credentials, sender and recipient only on the server, preferably via secrets files. Never expose them in public form metadata or browser bundles.
4. Resolve notification recipient from the authenticated, tenant-scoped form configuration. Never trust a recipient supplied by a public form submitter.
5. Sanitize submitted values and omit sensitive fields from notification bodies when configured.
6. Keep mail disabled by default until retention, anti-abuse, queueing and delivery tests pass.

The current transport supports explicit STARTTLS SMTP (typically port 587); implicit TLS on port 465 is not implemented. The SMTP server must advertise STARTTLS and present a valid certificate.

For Ousdal IT, plan separate forms for contact, support and courses, each with its own subject prefix and server-side recipient. The Ploform instance hostname is deployment-specific.
