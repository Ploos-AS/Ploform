-- Opaque sessions: never store plaintext bearer tokens.
CREATE TABLE sessions (
 id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
 user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
 token_digest bytea NOT NULL UNIQUE CHECK (octet_length(token_digest) = 32),
 created_at timestamptz NOT NULL DEFAULT now(),
 expires_at timestamptz NOT NULL,
 revoked_at timestamptz,
 CHECK (expires_at > created_at)
);
CREATE INDEX sessions_user_idx ON sessions (user_id);
CREATE INDEX sessions_expiry_idx ON sessions (expires_at);
-- Intentionally no grants to ploform_app. Authentication must use a
-- narrowly scoped credential or security-reviewed database function.
