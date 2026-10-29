-- Add optional encrypted payload storage without migrating existing plaintext.
-- This migration does not enable encrypted submission APIs.
ALTER TABLE submissions
  ADD COLUMN encryption_mode text NOT NULL DEFAULT 'standard'
    CHECK (encryption_mode IN ('standard', 'sensitive', 'confidential')),
  ADD COLUMN encrypted_payload jsonb,
  ADD CONSTRAINT submissions_payload_mode_check CHECK (
    (encryption_mode = 'standard' AND encrypted_payload IS NULL)
    OR
    (encryption_mode IN ('sensitive', 'confidential')
      AND encrypted_payload IS NOT NULL
      AND response_json = '{}'::jsonb)
  ),
  ADD CONSTRAINT submissions_envelope_shape_check CHECK (
    encrypted_payload IS NULL OR (
      jsonb_typeof(encrypted_payload) = 'object'
      AND encrypted_payload ?& ARRAY['version','algorithm','key_id','nonce','ciphertext']
      AND jsonb_typeof(encrypted_payload->'version') = 'number'
      AND jsonb_typeof(encrypted_payload->'algorithm') = 'string'
      AND jsonb_typeof(encrypted_payload->'key_id') = 'string'
      AND jsonb_typeof(encrypted_payload->'nonce') = 'string'
      AND jsonb_typeof(encrypted_payload->'ciphertext') = 'string'
    )
  );
COMMENT ON COLUMN submissions.encrypted_payload IS
  'Versioned authenticated-encryption envelope; never store key material here.';
COMMENT ON COLUMN submissions.response_json IS
  'Legacy plaintext responses for standard mode only; encrypted modes require empty object.';
