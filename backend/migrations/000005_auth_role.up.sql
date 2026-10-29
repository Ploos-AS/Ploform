-- Auth-specific role has no tenant data privileges and cannot bypass RLS.
DO $$
BEGIN
 IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = 'ploform_auth') THEN
  CREATE ROLE ploform_auth LOGIN;
 END IF;
END $$;
GRANT USAGE ON SCHEMA public TO ploform_auth;
GRANT SELECT, INSERT ON TABLE sessions TO ploform_auth;
GRANT UPDATE (revoked_at) ON TABLE sessions TO ploform_auth;
GRANT SELECT (user_id, organization_id, role) ON TABLE organization_memberships TO ploform_auth;
-- Session lookup also requires SELECT on the columns used in its predicate.
-- Never grant organization provisioning, table ownership, BYPASSRLS or SUPERUSER.
