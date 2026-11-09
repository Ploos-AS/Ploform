-- Regression for membership lifecycle. Everything rolls back.
BEGIN;

INSERT INTO organizations (id, slug, name)
VALUES ('10000000-0000-4000-8000-000000000001', 'status-test-org', 'Status test');
INSERT INTO users (id, email)
VALUES ('20000000-0000-4000-8000-000000000001', 'status-test@example.invalid');

INSERT INTO organization_memberships (organization_id, user_id, role)
VALUES ('10000000-0000-4000-8000-000000000001', '20000000-0000-4000-8000-000000000001', 'viewer');

DO $$
DECLARE membership_status text;
BEGIN
 SELECT status INTO membership_status FROM organization_memberships
 WHERE user_id = '20000000-0000-4000-8000-000000000001';
 IF membership_status IS DISTINCT FROM 'active' THEN
  RAISE EXCEPTION 'new membership must default to active';
 END IF;
END $$;

UPDATE organization_memberships SET status = 'suspended'
WHERE user_id = '20000000-0000-4000-8000-000000000001';

DO $$
BEGIN
 IF EXISTS (
   SELECT 1 FROM organization_memberships
   WHERE user_id = '20000000-0000-4000-8000-000000000001'
   AND organization_id = '10000000-0000-4000-8000-000000000001'
   AND status = 'active'
 ) THEN RAISE EXCEPTION 'suspended membership incorrectly active'; END IF;
END $$;

UPDATE organization_memberships SET status = 'revoked'
WHERE user_id = '20000000-0000-4000-8000-000000000001';

DO $$
BEGIN
 IF EXISTS (
   SELECT 1 FROM organization_memberships
   WHERE user_id = '20000000-0000-4000-8000-000000000001'
   AND status = 'active'
 ) THEN RAISE EXCEPTION 'revoked membership incorrectly active'; END IF;
 IF NOT EXISTS (
   SELECT 1 FROM information_schema.column_privileges
   WHERE table_schema = 'public'
   AND table_name = 'organization_memberships'
   AND column_name = 'status'
   AND grantee = 'ploform_auth'
   AND privilege_type = 'SELECT'
 ) THEN RAISE EXCEPTION 'auth role lacks SELECT on membership status'; END IF;
END $$;

ROLLBACK;
