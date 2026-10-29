-- Development bootstrap role only. Production credentials must be provisioned
-- securely by the operator, not committed to source control.
DO $$
BEGIN
 IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = 'ploform_app') THEN
  CREATE ROLE ploform_app LOGIN;
 END IF;
END $$;

GRANT USAGE ON SCHEMA public TO ploform_app;
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE workspaces, forms, form_versions, submissions TO ploform_app;
-- Organizations require a separate privileged provisioning workflow.
-- Never grant BYPASSRLS, table ownership or superuser to ploform_app.
