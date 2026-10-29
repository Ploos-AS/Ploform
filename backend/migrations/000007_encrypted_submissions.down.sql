ALTER TABLE submissions
  DROP CONSTRAINT submissions_envelope_shape_check,
  DROP CONSTRAINT submissions_payload_mode_check,
  DROP COLUMN encrypted_payload,
  DROP COLUMN encryption_mode;
