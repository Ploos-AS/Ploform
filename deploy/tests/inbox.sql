-- Inbox listing: tenant isolation, form separation, ordering and pagination.
-- Fixtures are temporary and every assertion runs under the restricted app role.
\set ON_ERROR_STOP on
BEGIN;
INSERT INTO organizations(slug,name) VALUES ('inbox-test-a','Inbox A') RETURNING id AS org_a \gset
INSERT INTO organizations(slug,name) VALUES ('inbox-test-b','Inbox B') RETURNING id AS org_b \gset
INSERT INTO workspaces(organization_id,slug,name) VALUES (:'org_a','inbox-ws','Inbox') RETURNING id AS ws_a \gset
INSERT INTO workspaces(organization_id,slug,name) VALUES (:'org_b','inbox-ws','Inbox') RETURNING id AS ws_b \gset
INSERT INTO forms(organization_id,workspace_id,slug,title) VALUES (:'org_a',:'ws_a','form-one','One') RETURNING id AS form_a \gset
INSERT INTO forms(organization_id,workspace_id,slug,title) VALUES (:'org_a',:'ws_a','form-two','Two') RETURNING id AS form_a2 \gset
INSERT INTO forms(organization_id,workspace_id,slug,title) VALUES (:'org_b',:'ws_b','form-one','Other') RETURNING id AS form_b \gset
INSERT INTO form_versions(organization_id,form_id,version,schema_json) VALUES (:'org_a',:'form_a',1,'{"fields":[]}'::jsonb) RETURNING id AS version_a \gset
INSERT INTO form_versions(organization_id,form_id,version,schema_json) VALUES (:'org_a',:'form_a2',1,'{"fields":[]}'::jsonb) RETURNING id AS version_a2 \gset
INSERT INTO form_versions(organization_id,form_id,version,schema_json) VALUES (:'org_b',:'form_b',1,'{"fields":[]}'::jsonb) RETURNING id AS version_b \gset
INSERT INTO submissions(organization_id,form_version_id,response_json,created_at) VALUES
 (:'org_a',:'version_a','{"message":"old"}',now()-interval '2 minutes'),
 (:'org_a',:'version_a','{"message":"new"}',now()-interval '1 minute'),
 (:'org_a',:'version_a2','{"message":"different form"}',now()),
 (:'org_b',:'version_b','{"message":"other tenant"}',now());
SET LOCAL ROLE ploform_app;
SELECT set_config('app.organization_id',:'org_a',true);
-- Mirrors the store query, including both tenant and form filters.
SELECT CASE WHEN count(*)=2 THEN 'true' ELSE 'false' END AS correct_count
FROM submissions s
JOIN form_versions v ON v.id=s.form_version_id AND v.organization_id=s.organization_id
JOIN forms f ON f.id=v.form_id AND f.organization_id=v.organization_id
WHERE s.organization_id=:'org_a' AND f.id=:'form_a' \gset
\if :correct_count
\else
\echo 'Inbox form filtering failed'
\quit 1
\endif
DO $$
BEGIN
 IF (SELECT count(*) FROM submissions) <> 3 THEN
  RAISE EXCEPTION 'tenant A can see submissions from tenant B';
 END IF;
END $$;
SELECT response_json->>'message' AS first_message
FROM submissions s
JOIN form_versions v ON v.id=s.form_version_id AND v.organization_id=s.organization_id
JOIN forms f ON f.id=v.form_id AND f.organization_id=v.organization_id
WHERE s.organization_id=:'org_a' AND f.id=:'form_a'
ORDER BY s.created_at DESC,s.id DESC LIMIT 1 OFFSET 0 \gset
SELECT CASE WHEN :'first_message'='new' THEN 'true' ELSE 'false' END AS correct_order \gset
\if :correct_order
\else
\echo 'Missing first inbox page'
\quit 1
\endif
DO $$
BEGIN
 IF (SELECT count(*) FROM submissions WHERE response_json->>'message'='other tenant')<>0 THEN
  RAISE EXCEPTION 'cross-tenant submission visible';
 END IF;
END $$;
SELECT set_config('app.organization_id',:'org_b',true);
DO $$
BEGIN
 IF (SELECT count(*) FROM submissions) <> 1 THEN
  RAISE EXCEPTION 'tenant B submission visibility incorrect';
 END IF;
END $$;
ROLLBACK;
\echo 'Transactional inbox isolation checks passed'
