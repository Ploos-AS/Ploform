-- Add an explicit lifecycle state without altering existing memberships.
-- Existing memberships remain active until an authorized administrator changes them.
ALTER TABLE organization_memberships
  ADD COLUMN status text NOT NULL DEFAULT 'active'
  CHECK (status IN ('active', 'suspended', 'revoked'));

CREATE INDEX organization_memberships_active_user_idx
  ON organization_memberships (user_id, organization_id)
  WHERE status = 'active';

-- The auth role must be able to read status to reject inactive memberships.
GRANT SELECT (status) ON TABLE organization_memberships TO ploform_auth;
