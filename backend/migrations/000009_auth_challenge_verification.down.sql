REVOKE SELECT (id, user_id, token_digest, expires_at, consumed_at, attempts) ON login_challenges FROM ploform_auth;
REVOKE UPDATE (consumed_at, attempts) ON login_challenges FROM ploform_auth;
