-- Verify real database session lifecycle and restricted auth role.
\set ON_ERROR_STOP on
BEGIN;
INSERT INTO users(email) VALUES ('auth-test-' || substr(gen_random_uuid()::text,1,12) || '@example.invalid') RETURNING id AS user_id \gset
INSERT INTO organizations(slug,name) VALUES ('auth-test-' || substr(gen_random_uuid()::text,1,12),'Auth test') RETURNING id AS org_id \gset
INSERT INTO organization_memberships(organization_id,user_id,role) VALUES (:'org_id',:'user_id','viewer');
INSERT INTO sessions(user_id,token_digest,expires_at) VALUES
 (:'user_id',digest('active-auth-test','sha256'),now()+interval '1 hour'),
 (:'user_id',digest('expired-auth-test','sha256'),now()-interval '1 hour'),
 (:'user_id',digest('revoked-auth-test','sha256'),now()+interval '1 hour');
UPDATE sessions SET revoked_at=now() WHERE token_digest=digest('revoked-auth-test','sha256');
SET LOCAL ROLE ploform_auth;
SELECT CASE WHEN count(*)=1 THEN 'true' ELSE 'false' END AS active_ok FROM sessions
 WHERE user_id=:'user_id' AND revoked_at IS NULL AND expires_at>now() \gset
\if :active_ok
\else
\echo 'Active session visibility failed'
\quit 1
\endif
SELECT CASE WHEN count(*)=0 THEN 'true' ELSE 'false' END AS expired_denied FROM sessions
 WHERE token_digest=digest('expired-auth-test','sha256') AND revoked_at IS NULL AND expires_at>now() \gset
\if :expired_denied
\else
\echo 'Expired session accepted'
\quit 1
\endif
SELECT CASE WHEN count(*)=0 THEN 'true' ELSE 'false' END AS revoked_denied FROM sessions
 WHERE token_digest=digest('revoked-auth-test','sha256') AND revoked_at IS NULL AND expires_at>now() \gset
\if :revoked_denied
\else
\echo 'Revoked session accepted'
\quit 1
\endif
SELECT CASE WHEN count(*)=1 THEN 'true' ELSE 'false' END AS membership_ok FROM organization_memberships
 WHERE user_id=:'user_id' AND organization_id=:'org_id' AND role='viewer' \gset
\if :membership_ok
\else
\echo 'Membership lookup failed'
\quit 1
\endif
ROLLBACK;
\echo 'Transactional auth session lifecycle checks passed'
