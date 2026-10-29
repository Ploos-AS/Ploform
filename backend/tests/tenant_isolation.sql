-- Run only against a disposable test database after migrations 000001 and 000002.
-- Must connect as the migration owner/superuser. All changes roll back.
\set ON_ERROR_STOP on
BEGIN;
INSERT INTO organizations (id, slug, name) VALUES
 ('11111111-1111-4111-8111-111111111111','tenant-alpha','Alpha'),
 ('22222222-2222-4222-8222-222222222222','tenant-beta','Beta');
INSERT INTO workspaces (id, organization_id, slug, name) VALUES
 ('aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa','11111111-1111-4111-8111-111111111111','alpha-ws','Alpha workspace'),
 ('bbbbbbbb-bbbb-4bbb-8bbb-bbbbbbbbbbbb','22222222-2222-4222-8222-222222222222','beta-ws','Beta workspace');

SET LOCAL ROLE ploform_app;
SELECT set_config('app.organization_id','11111111-1111-4111-8111-111111111111',true);

DO $$
DECLARE visible_count int;
BEGIN
 SELECT count(*) INTO visible_count FROM workspaces;
 IF visible_count <> 1 THEN RAISE EXCEPTION 'RLS read isolation failed: expected 1, got %', visible_count; END IF;
 IF EXISTS (SELECT 1 FROM workspaces WHERE id='bbbbbbbb-bbbb-4bbb-8bbb-bbbbbbbbbbbb') THEN
  RAISE EXCEPTION 'Tenant B workspace leaked to tenant A';
 END IF;
END $$;

-- A forged cross-tenant insert must be rejected by RLS WITH CHECK.
DO $$
BEGIN
 BEGIN
  INSERT INTO workspaces (organization_id,slug,name)
  VALUES ('22222222-2222-4222-8222-222222222222','forged','Forged');
  RAISE EXCEPTION 'Cross-tenant insert was unexpectedly accepted';
 EXCEPTION WHEN insufficient_privilege THEN
  NULL;
 END;
END $$;

-- A cross-tenant update must affect zero rows.
DO $$
DECLARE affected int;
BEGIN
 UPDATE workspaces SET name='Stolen' WHERE id='bbbbbbbb-bbbb-4bbb-8bbb-bbbbbbbbbbbb';
 GET DIAGNOSTICS affected = ROW_COUNT;
 IF affected <> 0 THEN RAISE EXCEPTION 'Cross-tenant update affected % rows', affected; END IF;
END $$;

-- A cross-tenant delete must affect zero rows.
DO $$
DECLARE affected int;
BEGIN
 DELETE FROM workspaces WHERE id='bbbbbbbb-bbbb-4bbb-8bbb-bbbbbbbbbbbb';
 GET DIAGNOSTICS affected = ROW_COUNT;
 IF affected <> 0 THEN RAISE EXCEPTION 'Cross-tenant delete affected % rows', affected; END IF;
END $$;

RESET ROLE;
ROLLBACK;
