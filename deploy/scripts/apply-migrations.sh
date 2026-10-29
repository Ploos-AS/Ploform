#!/bin/sh
# Development-only migration runner. SQL file + history record are atomic.
set -eu
psql -X -v ON_ERROR_STOP=1 -c '
CREATE TABLE IF NOT EXISTS ploform_schema_migrations (
  name text PRIMARY KEY,
  applied_at timestamptz NOT NULL DEFAULT now()
);'
for migration in \
  000001_initial \
  000002_app_role \
  000003_identity \
  000004_sessions \
  000005_auth_role \
  000006_login_challenges \
  000007_encrypted_submissions \
  000008_webhook_outbox
do
  file="/migrations/${migration}.up.sql"
  [ -r "$file" ] || { echo "Missing migration: $file" >&2; exit 1; }
  applied="$(psql -X -qAt -v ON_ERROR_STOP=1 -c "SELECT 1 FROM ploform_schema_migrations WHERE name = '$migration'")"
  if [ "$applied" = 1 ]; then
    echo "Already applied: $migration"
    continue
  fi
  echo "Applying $migration"
  psql -X -v ON_ERROR_STOP=1 -1 \
    -f "$file" \
    -c "INSERT INTO ploform_schema_migrations(name) VALUES ('$migration')"
done
echo "Development schema migrations complete"
