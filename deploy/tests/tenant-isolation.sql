-- Transactional tenant isolation integration test. No persistent test data.
\set ON_ERROR_STOP on
BEGIN;
INSERT INTO organizations (slug, name)
VALUES ('integration-tenant-a', 'Integration Tenant A')
RETURNING id AS org_a \gset
INSERT INTO organizations (slug, name)
VALUES ('integration-tenant-b', 'Integration Tenant B')
RETURNING id AS org_b \gset
INSERT INTO workspaces (organization_id, slug, name)
VALUES (:'org_a', 'integration-workspace', 'Test Workspace')
RETURNING id AS workspace_a \gset
INSERT INTO workspaces (organization_id, slug, name)
VALUES (:'org_b', 'integration-workspace', 'Test Workspace')
RETURNING id AS workspace_b \gset

-- Execute tenant queries as the actual restricted application role.
SET LOCAL ROLE ploform_app;
SELECT set_config('app.organization_id', :'org_a', true);
INSERT INTO forms (organization_id, workspace_id, slug, title)
VALUES (:'org_a', :'workspace_a', 'integration-form', 'Integration Form')
RETURNING id AS form_a \gset
INSERT INTO form_versions (organization_id, form_id, version, schema_json)
VALUES (:'org_a', :'form_a', 1, '{}'::jsonb);
INSERT INTO form_versions (organization_id, form_id, version, schema_json)
VALUES (:'org_a', :'form_a', 2, '{}'::jsonb);

DO $$
BEGIN
  IF (SELECT count(*) FROM forms) <> 1 THEN
    RAISE EXCEPTION 'tenant A cannot read its own form';
  END IF;
  IF (SELECT count(*) FROM form_versions) <> 2 THEN
    RAISE EXCEPTION 'tenant A version read/count failed';
  END IF;
  IF (SELECT max(version) FROM form_versions) <> 2 THEN
    RAISE EXCEPTION 'version ordering failed';
  END IF;
END $$;

SELECT set_config('app.organization_id', :'org_b', true);
DO $$
BEGIN
  IF EXISTS (SELECT 1 FROM forms) THEN
    RAISE EXCEPTION 'RLS LEAK: tenant B can see tenant A form';
  END IF;
  IF EXISTS (SELECT 1 FROM form_versions) THEN
    RAISE EXCEPTION 'RLS LEAK: tenant B can see tenant A versions';
  END IF;
END $$;
INSERT INTO forms (organization_id, workspace_id, slug, title)
VALUES (:'org_b', :'workspace_b', 'integration-form', 'Integration Form B');
DO $$
BEGIN
  IF (SELECT count(*) FROM forms) <> 1 THEN
    RAISE EXCEPTION 'tenant B own form visibility failed';
  END IF;
END $$;

SELECT set_config('app.organization_id', :'org_a', true);
DO $$
BEGIN
  IF (SELECT count(*) FROM forms) <> 1 THEN
    RAISE EXCEPTION 'tenant A can see another tenant form';
  END IF;
END $$;
ROLLBACK;
\echo 'Transactional tenant RLS and form version checks passed'
