# SMTP worker (development)

`notify.Worker` polls the PostgreSQL outbox and invokes the existing authenticated STARTTLS SMTP transport.

## Safety and deployment status

- Worker is **not started by the API**. This commit does not send any messages.
- The worker requires a database, valid SMTP configuration, and an explicitly approved recipient. It does not accept recipient addresses from public submissions.
- SMTP secrets must be supplied only through server-side configuration.
- The worker does not log submitted content, recipients, or credentials.
- Failed SMTP sends are retried by the existing outbox mechanism.

## Remaining before end-to-end delivery

1. Review the existing migrations and connect notification enqueue to submission storage **in the same transaction**.
2. Add tenant-scoped, server-managed notification configuration and safe per-form routing.
3. Add a dedicated worker executable with secret-file support and graceful shutdown.
4. Test PostgreSQL migrations, transactional rollback, retry, concurrency and SMTP integration.
5. Review abuse protection, retention, privacy disclosures, and duplicate email behavior.

**No public site should enable email notifications until these checks pass.**
