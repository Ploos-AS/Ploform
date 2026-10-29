-- Run as migration owner/superuser in an isolated CI database.
BEGIN;

DO $$
DECLARE
  org_id uuid := gen_random_uuid();
  user_id uuid := gen_random_uuid();
  digest bytea := decode(repeat('ab', 32), 'hex');
BEGIN
  INSERT INTO organizations (id, slug, name)
  VALUES (org_id, 'auth-role-ci', 'Auth role test');
  INSERT INTO users (id, email)
  VALUES (user_id, 'auth-role-test@example.invalid');
  INSERT INTO organization_memberships (organization_id, user_id, role)
  VALUES (org_id, user_id, 'viewer');
  INSERT INTO sessions (user_id, token_digest, expires_at)
  VALUES (user_id, digest, now() + interval '1 hour');
END $$;

SET LOCAL ROLE ploform_auth;

DO $$
DECLARE
  session_count integer;
  membership_count integer;
BEGIN
  SELECT count(*) INTO session_count FROM sessions;
  IF session_count != 1 THEN
    RAISE EXCEPTION 'auth role cannot read sessions';
  END IF;

  SELECT count(*) INTO membership_count FROM organization_memberships;
  IF membership_count != 1 THEN
    RAISE EXCEPTION 'auth role cannot read memberships';
  END IF;

  -- A session service must not read tenant data directly.
  BEGIN
    PERFORM 1 FROM workspaces;
    RAISE EXCEPTION 'auth role unexpectedly accessed workspaces';
  EXCEPTION WHEN insufficient_privilege THEN
    NULL;
  END;

  -- An auth role must not create users.
  BEGIN
    INSERT INTO users (email) VALUES ('forbidden@example.invalid');
    RAISE EXCEPTION 'auth role unexpectedly inserted a user';
  EXCEPTION WHEN insufficient_privilege THEN
    NULL;
  END;
END $$;

RESET ROLE;
ROLLBACK;
