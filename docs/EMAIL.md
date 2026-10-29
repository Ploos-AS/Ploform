# Outbound email design (proposal)

Self-hosted deployments should be able to choose their own outbound
provider. Start with generic SMTP submission over authenticated TLS,
including Proton Mail SMTP submission where the account supports it.
Later adapters may include SES and transactional email APIs.

Keep transport credentials out of source control and logs. Require explicit
TLS verification, sender/domain validation, rate limits, retry queues,
delivery observability and safe failure handling.

System messages (login challenges) and form notifications must use separate
templates and rate limits. By default, notification messages contain only
a generic notice and a link to authenticated Ploform access, never submitted
answers or attachments.

SMTP is not synonymous with end-to-end encryption. Optional OpenPGP
encrypted notifications require recipient key verification, secure key
rotation and a warning about exposed email metadata.

No outbound email adapter is implemented yet.
