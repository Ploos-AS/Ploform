-- Contact-form persistence integration test using the restricted app role.
-- All fixtures are rolled back, including deliberate constraint failures.
\set ON_ERROR_STOP on
BEGIN;
INSERT INTO organizations (slug, name)
VALUES ('integration-contact-org', 'Contact Integration')
RETURNING id AS org_id \gset
INSERT INTO workspaces (organization_id, slug, name)
VALUES (:'org_id', 'contact-workspace', 'Contact Workspace')
RETURNING id AS workspace_id \gset
SET LOCAL ROLE ploform_app;
SELECT set_config('app.organization_id', :'org_id', true);

-- Mirrors the two INSERTs in Store.CreateContactForm within one transaction.
SAVEPOINT contact_creation;
INSERT INTO forms (organization_id, workspace_id, slug, title)
VALUES (:'org_id', :'workspace_id', 'contact', 'Kontakt oss')
RETURNING id AS form_id \gset
INSERT INTO form_versions (organization_id, form_id, version, schema_json)
VALUES (:'org_id', :'form_id', 1,
 '{"fields":[{"id":"name","type":"text","label":"Navn","required":true},{"id":"email","type":"email","label":"E-post","required":true},{"id":"subject","type":"text","label":"Emne"},{"id":"message","type":"textarea","label":"Melding","required":true}]}'::jsonb);
-- Assert persisted schema fields, including required flags.
DO $$
DECLARE fields jsonb;
BEGIN
 SELECT schema_json->'fields' INTO fields FROM form_versions WHERE version=1;
 IF jsonb_array_length(fields) <> 4
    OR fields->0->>'id' <> 'name'
    OR fields->1->>'id' <> 'email'
    OR fields->2->>'id' <> 'subject'
    OR fields->3->>'id' <> 'message'
    OR fields->0->>'required' <> 'true'
    OR fields->1->>'required' <> 'true'
    OR fields->2 ? 'required'
    OR fields->3->>'required' <> 'true'
 THEN RAISE EXCEPTION 'contact schema differs from expected fields'; END IF;
END $$;

-- Simulate an error after inserting a form, and verify savepoint rollback
-- leaves no orphan draft.
SAVEPOINT failed_contact;
INSERT INTO forms (organization_id, workspace_id, slug, title)
VALUES (:'org_id', :'workspace_id', 'failed-contact', 'Should Roll Back');
ROLLBACK TO SAVEPOINT failed_contact;
DO $$
BEGIN
 IF EXISTS (SELECT 1 FROM forms WHERE slug='failed-contact') THEN
  RAISE EXCEPTION 'failed contact form was not rolled back';
 END IF;
END $$;
-- Publication: only a versioned draft can transition to published.
INSERT INTO forms (organization_id, workspace_id, slug, title)
VALUES (:'org_id', :'workspace_id', 'unversioned', 'Unversioned');
UPDATE forms SET status='published'
WHERE organization_id=:'org_id' AND id=:'form_id'
  AND status='draft'
  AND EXISTS (SELECT 1 FROM form_versions
              WHERE organization_id=:'org_id' AND form_id=:'form_id');
DO $
BEGIN
 IF (SELECT status FROM forms WHERE slug='contact') <> 'published' THEN
  RAISE EXCEPTION 'versioned contact form was not published';
 END IF;
 IF (SELECT status FROM forms WHERE slug='unversioned') <> 'draft' THEN
  RAISE EXCEPTION 'unversioned form was published';
 END IF;
END $;
-- Repeated publication must not change the already-published form.
UPDATE forms SET status='published'
WHERE organization_id=:'org_id' AND id=:'form_id' AND status='draft';
DO $
BEGIN
 IF FOUND THEN RAISE EXCEPTION 'published form accepted duplicate transition'; END IF;
END $;

ROLLBACK;
\echo 'Transactional contact form persistence and publication checks passed'
