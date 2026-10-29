REVOKE ALL PRIVILEGES ON TABLE sessions, organization_memberships FROM ploform_auth;
REVOKE USAGE ON SCHEMA public FROM ploform_auth;
-- The role is retained to avoid affecting other databases.
