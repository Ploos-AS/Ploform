#!/bin/sh
# Development-only bootstrap. Fail closed on any migration error.
set -eu
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
  echo "Applying $migration"
  psql -X -v ON_ERROR_STOP=1 -1 -f "$file"
done
echo "Development schema migrations complete"
