-- Identity and membership schema; no login endpoints are enabled.
CREATE TABLE users (
 id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
 email text NOT NULL UNIQUE CHECK (length(trim(email)) > 3),
 created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE organization_memberships (
 organization_id uuid NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
 user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
 role text NOT NULL CHECK (role IN ('owner','admin','editor','viewer')),
 created_at timestamptz NOT NULL DEFAULT now(),
 PRIMARY KEY (organization_id, user_id)
);
CREATE INDEX organization_memberships_user_idx ON organization_memberships (user_id);

-- Membership lookup is intentionally not granted to ploform_app yet.
-- A future authentication service must resolve membership securely and
-- pass only authorized organization IDs to tenant-scoped transactions.
