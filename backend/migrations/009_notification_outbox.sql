-- Transactional email notification outbox.
-- Apply using the repository's existing migration mechanism before enabling workers.
CREATE TABLE IF NOT EXISTS notification_outbox (
 id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
 organization_id uuid NOT NULL,
 submission_id uuid NOT NULL,
 recipient text NOT NULL CHECK (length(recipient) BETWEEN 3 AND 320),
 subject text NOT NULL CHECK (length(subject) BETWEEN 1 AND 500),
 body text NOT NULL,
 status text NOT NULL DEFAULT 'pending' CHECK (status IN ('pending','sending','retry','sent','failed')),
 attempts integer NOT NULL DEFAULT 0 CHECK (attempts >= 0),
 next_attempt_at timestamptz NOT NULL DEFAULT now(),
 lease_until timestamptz,
 sent_at timestamptz,
 created_at timestamptz NOT NULL DEFAULT now(),
 updated_at timestamptz NOT NULL DEFAULT now(),
 UNIQUE (organization_id,submission_id)
);
CREATE INDEX IF NOT EXISTS notification_outbox_due_idx
 ON notification_outbox(next_attempt_at,id)
 WHERE status IN ('pending','retry','sending');
