REVOKE ALL PRIVILEGES ON TABLE workspaces, forms, form_versions, submissions FROM ploform_app;
REVOKE USAGE ON SCHEMA public FROM ploform_app;
-- Role intentionally retained; it may be used by other databases.
