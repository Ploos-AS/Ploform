-- The authentication role may verify and consume challenges, but cannot issue them.
-- Challenge issuance remains disabled pending email delivery and abuse controls.
GRANT SELECT (id, user_id, token_digest, expires_at, consumed_at, attempts) ON login_challenges TO ploform_auth;
GRANT UPDATE (consumed_at, attempts) ON login_challenges TO ploform_auth;
