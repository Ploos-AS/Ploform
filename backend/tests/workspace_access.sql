-- End-to-end database regression for membership, sessions and tenant RLS.
-- Run only on a disposable database as migration owner.
BEGIN;

CREATE TEMP TABLE test_ids (label text PRIMARY KEY, id uuid NOT NULL);
INSERT INTO organizations (slug, name) VALUES ('workspace-access-a', 'Workspace A') RETURNING id \gset
INSERT INTO test_ids VALUES ('org_a', :'id');
INSERT INTO organizations (slug, name) VALUES ('workspace-access-b', 'Workspace B') RETURNING id \gset
INSERT INTO test_ids VALUES ('org_b', :'id');
INSERT INTO users (email) VALUES ('workspace-access@example.invalid') RETURNING id \gset
INSERT INTO test_ids VALUES ('user', :'id');

INSERT INTO organization_memberships (organization_id, user_id, role)
SELECT (SELECT id FROM test_ids WHERE label = 'org_a'),
       (SELECT id FROM test_ids WHERE label = 'user'), 'viewer';

INSERT INTO sessions (user_id, token_digest, expires_at)
SELECT id, decode(repeat('cd', 32), 'hex'), now() + interval '1 hour'
FROM test_ids WHERE label = 'user';

INSERT INTO workspaces (organization_id, slug, name)
SELECT id, 'visible-workspace', 'Visible' FROM test_ids WHERE label = 'org_a';
INSERT INTO workspaces (organization_id, slug, name)
SELECT id, 'hidden-workspace', 'Hidden' FROM test_ids WHERE label = 'org_b';

SET LOCAL ROLE ploform_auth;
DO $$
DECLARE n integer;
BEGIN
 SELECT count(*) INTO n
 FROM sessions s
 JOIN organization_memberships m ON m.user_id = s.user_id
 WHERE s.token_digest = decode(repeat('cd', 32), 'hex')
   AND s.revoked_at IS NULL AND s.expires_at > now()
   AND m.organization_id = (SELECT id FROM test_ids WHERE label = 'org_a');
 IF n != 1 THEN RAISE EXCEPTION 'valid session/membership not resolved'; END IF;

 SELECT count(*) INTO n
 FROM sessions s
 JOIN organization_memberships m ON m.user_id = s.user_id
 WHERE s.token_digest = decode(repeat('cd', 32), 'hex')
   AND m.organization_id = (SELECT id FROM test_ids WHERE label = 'org_b');
 IF n != 0 THEN RAISE EXCEPTION 'cross-organization membership exposed'; END IF;
END $$;
RESET ROLE;

SET LOCAL ROLE ploform_app;
SELECT set_config('app.organization_id', (SELECT id::text FROM test_ids WHERE label = 'org_a'), true);
DO $$
DECLARE n integer;
BEGIN
 SELECT count(*) INTO n FROM workspaces;
 IF n != 1 THEN RAISE EXCEPTION 'tenant A expected 1 workspace, found %', n; END IF;
 SELECT count(*) INTO n FROM workspaces WHERE slug = 'hidden-workspace';
 IF n != 0 THEN RAISE EXCEPTION 'tenant B workspace leaked'; END IF;
END $$;
RESET ROLE;

ROLLBACK;
