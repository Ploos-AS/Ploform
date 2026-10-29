#!/bin/sh
# Run on an operator machine with Docker Compose installed.
# Internal-only deployment verification; no public exposure.
set -eu
cd "$(dirname "$0")/../.."
compose="deploy/compose.proxy.dev.yaml"
: "${PLOFORM_DEV_POSTGRES_PASSWORD:?Set PLOFORM_DEV_POSTGRES_PASSWORD}"
command -v docker >/dev/null || { echo "Docker required" >&2; exit 1; }
for file in deploy/secrets/app_password.txt deploy/secrets/auth_password.txt; do
 [ -s "$file" ] || { echo "Missing $file" >&2; exit 1; }
done
docker compose -f "$compose" config --quiet
docker compose -f "$compose" up --build -d --wait
[ "$(docker compose -f "$compose" ps -q api)" != "" ]
# HTTP check from inside Caddy's isolated network (no published ports).
docker compose -f "$compose" exec -T caddy wget -qO- http://api:8080/healthz | grep -q '"status":"ok"'
docker compose -f "$compose" exec -T caddy wget -qO- http://localhost/healthz | grep -q '"status":"ok"'
# Assert restricted roles exist and cannot bypass row-level security.
docker compose -f "$compose" exec -T postgres psql -X -U ploform_dev -d ploform -v ON_ERROR_STOP=1 -At <<'SQL' | grep -qx 1
SELECT CASE WHEN count(*) = 2 THEN 1 ELSE 0 END
FROM pg_roles
WHERE rolname IN ('ploform_app', 'ploform_auth')
  AND NOT rolbypassrls AND NOT rolsuper;
SQL
# Confirm the API reports readiness (database PingContext).
docker compose -f "$compose" exec -T caddy wget -qO- http://api:8080/readyz >/dev/null
echo "Internal API and Caddy smoke checks passed"
