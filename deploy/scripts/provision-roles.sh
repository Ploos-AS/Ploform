#!/bin/sh
# Development-only role password provisioning. Never echo secrets.
set -eu
[ -s /run/secrets/ploform_app_password ] || exit 1
[ -s /run/secrets/ploform_auth_password ] || exit 1
app_password="$(cat /run/secrets/ploform_app_password)"
auth_password="$(cat /run/secrets/ploform_auth_password)"
[ "${#app_password}" -ge 24 ] || { echo "Application password too short" >&2; exit 1; }
[ "${#auth_password}" -ge 24 ] || { echo "Auth password too short" >&2; exit 1; }
[ "$app_password" != "$auth_password" ] || { echo "Role passwords must differ" >&2; exit 1; }
# psql variables are quoted as SQL literals; do not interpolate into SQL text.
psql -X -v ON_ERROR_STOP=1 -v app_password="$app_password" -v auth_password="$auth_password" <<'SQL'
SELECT format('ALTER ROLE ploform_app PASSWORD %L', :'app_password') \gexec
SELECT format('ALTER ROLE ploform_auth PASSWORD %L', :'auth_password') \gexec
SQL
unset app_password auth_password
echo "Development runtime roles provisioned"
