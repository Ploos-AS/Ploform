# Passwordless login design (not enabled)

The schema in migration 000006 introduces short-lived, single-use login
challenges. There are no public issuance or redemption endpoints.

Security requirements before activation:

- Prove ownership of an email address through a trusted delivery channel.
- Return indistinguishable responses for known and unknown accounts.
- Store challenge digests, never plaintext secrets.
- Limit lifetime to 15 minutes, enforce a strict attempt budget and
  atomically consume a challenge before issuing a session.
- Apply IP and account-based rate limits and delivery throttling.
- Never put session bearer tokens in URLs.
- Do not let callers choose user IDs, organization IDs or roles.
- Validate the external HTTPS origin, CSRF protections and reverse proxy.
- Audit challenge replay, expiry, revocation and cross-tenant access.

The new table intentionally has no grants for runtime database roles.
