# Secure Forms security design (proposal, not implemented)

Ploform will support three explicitly labeled modes:

- **Standard:** HTTPS in transit; database encryption at rest depends on deployment.
- **Sensitive:** application-level authenticated encryption of submission payloads
  with per-organization keys managed outside PostgreSQL. The API necessarily
  handles plaintext during submission and authorized reading.
- **Confidential:** browser-side encryption to a pinned recipient public key,
  with only ciphertext delivered to the API. The server still sees metadata
  and can compromise browser-delivered JavaScript; do not claim unconditional
  zero knowledge.

## Threat model

Protect submission contents from a leaked database backup, unauthorized
cross-tenant reads and accidental email disclosure. Separate threats include
compromised application servers, malicious client scripts, stolen recipient
private keys, traffic analysis, malicious form owners and browser extensions.
Encryption does not replace access control, anti-abuse or retention policies.

## Implementation rules

1. Never invent cryptographic primitives. Select an established, reviewed
   envelope encryption construction and versioned message format.
2. Use independent encryption contexts per organization; bind ciphertext
   authentication to organization ID, form ID, version and submission ID.
3. Keep private keys and master keys outside the database. Document key
   recovery, rotation, revocation and irrecoverable key loss.
4. Encrypt file attachments as well as structured answers; enforce size limits
   and malware-handling policy before enabling uploads.
5. Never put plaintext answers, access tokens or private keys in emails,
   URLs, application logs, analytics or webhook payloads.
6. For confidential mode, clearly communicate that the server can replace
   the delivered JavaScript and that browser security is in scope.
7. Minimize metadata and retention; support deletion and export without
   leaking cross-tenant data.
8. Require explicit tests for ciphertext tampering, key mismatch, tenant
   isolation, replay, rotation, backup restore and authorization.
9. Keep encrypted-submission endpoints disabled until security review and
   independent testing.

## Email

Send notification-only emails by default. Optional OpenPGP-encrypted delivery
may be added later, with recipient key verification and a warning that
subject lines, routing and other mail metadata are not end-to-end encrypted.

## Status

This is a design document, not a claim that Ploform currently encrypts
submissions at application or browser level.

## Implemented primitive (not connected to storage)

`backend/internal/secureforms/envelope.go` implements AES-256-GCM envelope encryption with a fresh random 96-bit nonce per operation. The authenticated additional data binds organization, form, form version and submission identifiers. Envelope metadata includes a version, algorithm identifier and key ID. The caller supplies a 32-byte key; there is **no key manager, persistent ciphertext storage, HTTP integration or browser-side encryption** yet. Keys must never be generated from passwords or stored in the submission database. Key ID authenticity and key-to-organization binding must be enforced by the future key manager, not trusted from unverified envelope metadata. Unit tests have been written but not executed.

## Development-only key resolver

`backend/internal/secureforms/keyring.go` implements an in-memory resolver that binds a globally unique key ID to one organization. It copies key bytes on registration and retrieval and rejects cross-organization lookups. This is **not production key management**: keys disappear on restart, are resident in process memory, have no rotation or recovery, and the resolver does not authenticate callers. Do not expose it through HTTP or use it to protect real submissions. Before deployment, implement an external KMS/HSM or audited envelope key management system with authorization, per-organization key isolation, backup/recovery, audit logging, key rotation and documented failure behavior. Tests have been written but not executed.
