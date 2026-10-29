-- One-time, short-lived login challenges. No HTTP login flow is enabled.
CREATE TABLE login_challenges (
 id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
 user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
 token_digest bytea NOT NULL UNIQUE CHECK (octet_length(token_digest) = 32),
 created_at timestamptz NOT NULL DEFAULT now(),
 expires_at timestamptz NOT NULL,
 consumed_at timestamptz,
 attempts integer NOT NULL DEFAULT 0 CHECK (attempts >= 0 AND attempts <= 5),
 CHECK (expires_at > created_at),
 CHECK (expires_at <= created_at + interval '15 minutes')
);
CREATE INDEX login_challenges_user_idx ON login_challenges (user_id);
CREATE INDEX login_challenges_expiry_idx ON login_challenges (expires_at);
-- No grants to ploform_app or ploform_auth: challenge verification and
-- issuance need a separate security review before exposing credentials.
