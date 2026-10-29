-- Run in a disposable PostgreSQL database after migrations 000001-000007.
BEGIN;
DO $$
DECLARE
  org_id uuid := gen_random_uuid();
  workspace_id uuid := gen_random_uuid();
  form_id uuid := gen_random_uuid();
  version_id uuid := gen_random_uuid();
  payload jsonb := '{"version":1,"algorithm":"AES-256-GCM","key_id":"key-1","nonce":"nonce","ciphertext":"ciphertext"}';
BEGIN
  INSERT INTO organizations(id, slug, name) VALUES(org_id, 'encrypted-test', 'Encrypted test');
  INSERT INTO workspaces(id, organization_id, slug, name) VALUES(workspace_id, org_id, 'secure', 'Secure');
  INSERT INTO forms(id, organization_id, workspace_id, slug, title) VALUES(form_id, org_id, workspace_id, 'secure', 'Secure');
  INSERT INTO form_versions(id, organization_id, form_id, version, schema_json)
    VALUES(version_id, org_id, form_id, 1, '{}'::jsonb);
  INSERT INTO submissions(organization_id, form_version_id, response_json, encryption_mode, encrypted_payload)
    VALUES(org_id, version_id, '{}'::jsonb, 'sensitive', payload);
  BEGIN
    INSERT INTO submissions(organization_id, form_version_id, response_json, encryption_mode, encrypted_payload)
      VALUES(org_id, version_id, '{"secret":"plaintext"}'::jsonb, 'sensitive', payload);
    RAISE EXCEPTION 'plaintext accepted in sensitive mode';
  EXCEPTION WHEN check_violation THEN NULL;
  END;
  BEGIN
    INSERT INTO submissions(organization_id, form_version_id, response_json, encryption_mode)
      VALUES(org_id, version_id, '{}'::jsonb, 'confidential');
    RAISE EXCEPTION 'missing envelope accepted';
  EXCEPTION WHEN check_violation THEN NULL;
  END;
END $$;
ROLLBACK;
