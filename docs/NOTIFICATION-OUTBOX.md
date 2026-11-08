# SMTP notification outbox — development foundation

The SQL outbox and Go worker primitives support atomic enqueueing, leasing, bounded retries, and marking sent messages.

## Integration still required

- Verify and adapt the SQL migration filename/location to the repository's migration runner.
- Enqueue in the **same database transaction** as the persisted form submission. Do not enqueue separately.
- Determine the tenant-scoped recipient from authenticated server configuration, not public input.
- Run a background worker that calls `Outbox.ProcessOne` with `SMTPConfig.Send`; configure secrets on the server.
- Verify migration, integration tests, concurrency, crash recovery and privacy controls before deployment.
- Review retention for failed jobs and sensitive message bodies.

## Delivery semantics

SMTP cannot provide exactly-once delivery. If SMTP accepts a message and the worker crashes before recording success, a retry can deliver a duplicate. The unique enqueue key prevents duplicate jobs, **not duplicate emails**. Use a stable Message-ID and downstream deduplication where supported.

The outbox is **not yet connected to form submission**. No email is currently sent automatically.
