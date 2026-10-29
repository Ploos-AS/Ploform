-- Exercise challenge consumption with the same restricted auth role as the API.
\set ON_ERROR_STOP on
BEGIN;
INSERT INTO users(email) VALUES ('challenge-' || substr(gen_random_uuid()::text,1,12) || '@example.invalid') RETURNING id AS user_id \gset
INSERT INTO login_challenges(user_id,token_digest,expires_at,attempts) VALUES
 (:'user_id',digest('valid-challenge','sha256'),now()+interval '10 minutes',0),
 (:'user_id',digest('expired-challenge','sha256'),now()+interval '1 second',0),
 (:'user_id',digest('locked-challenge','sha256'),now()+interval '10 minutes',5),
 (:'user_id',digest('invalidated-challenge','sha256'),now()+interval '10 minutes',0);
UPDATE login_challenges SET consumed_at=clock_timestamp()
 WHERE token_digest=digest('invalidated-challenge','sha256');
SELECT pg_sleep(2);
SET LOCAL ROLE ploform_auth;
-- A valid token is consumed once, returning the expected user.
WITH consumed AS (
 UPDATE login_challenges SET consumed_at=clock_timestamp()
 WHERE token_digest=digest('valid-challenge','sha256')
 AND consumed_at IS NULL AND expires_at>clock_timestamp() AND attempts<5
 RETURNING user_id
)
SELECT (count(*)=1 AND bool_and(user_id=:'user_id'::uuid))::int AS valid_once FROM consumed \gset
\if :valid_once
\else
SELECT 1/0;
\endif
-- Replaying the same token must not consume a second row.
WITH consumed AS (
 UPDATE login_challenges SET consumed_at=clock_timestamp()
 WHERE token_digest=digest('valid-challenge','sha256')
 AND consumed_at IS NULL AND expires_at>clock_timestamp() AND attempts<5
 RETURNING user_id
)
SELECT (count(*)=0)::int AS replay_denied FROM consumed \gset
\if :replay_denied
\else
SELECT 1/0;
\endif
WITH consumed AS (
 UPDATE login_challenges SET consumed_at=clock_timestamp()
 WHERE token_digest=digest('expired-challenge','sha256')
 AND consumed_at IS NULL AND expires_at>clock_timestamp() AND attempts<5
 RETURNING user_id
)
SELECT (count(*)=0)::int AS expired_denied FROM consumed \gset
\if :expired_denied
\else
SELECT 1/0;
\endif
WITH consumed AS (
 UPDATE login_challenges SET consumed_at=clock_timestamp()
 WHERE token_digest=digest('locked-challenge','sha256')
 AND consumed_at IS NULL AND expires_at>clock_timestamp() AND attempts<5
 RETURNING user_id
)
SELECT (count(*)=0)::int AS locked_denied FROM consumed \gset
\if :locked_denied
\else
SELECT 1/0;
\endif
WITH consumed AS (
 UPDATE login_challenges SET consumed_at=clock_timestamp()
 WHERE token_digest=digest('unknown-challenge','sha256')
 AND consumed_at IS NULL AND expires_at>clock_timestamp() AND attempts<5
 RETURNING user_id
)
SELECT (count(*)=0)::int AS unknown_denied FROM consumed \gset
\if :unknown_denied
\else
SELECT 1/0;
\endif
WITH consumed AS (
 UPDATE login_challenges SET consumed_at=clock_timestamp()
 WHERE token_digest=digest('invalidated-challenge','sha256')
 AND consumed_at IS NULL AND expires_at>clock_timestamp() AND attempts<5
 RETURNING user_id
)
SELECT (count(*)=0)::int AS consumed_denied FROM consumed \gset
\if :consumed_denied
\else
SELECT 1/0;
\endif
ROLLBACK;
\echo 'Transactional one-time challenge checks passed'
