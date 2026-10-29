#!/bin/sh
# Run on an operator machine with Docker Compose and curl installed.
# Internal-only deployment verification; no public exposure.
set -eu
cd "$(dirname "$0")/../.."
compose="deploy/compose.proxy.dev.yaml"
: "${PLOFORM_DEV_POSTGRES_PASSWORD:?Set PLOFORM_DEV_POSTGRES_PASSWORD}"
command -v docker >/dev/null || { echo "Docker required" >&2; exit 1; }
for file in deploy/secrets/app_password.txt deploy/secrets/auth_password.txt; do
 [ -s "$file" ] || { echo "Missing $file" >&2; exit 1; }
done
docker compose -f "$compose" up --build -d
[ "$(docker compose -f "$compose" ps -q api)" != "" ]
# HTTP check from inside Caddy's isolated network (no published ports).
docker compose -f "$compose" exec -T caddy wget -qO- http://api:8080/healthz | grep -q '"status":"ok"'
docker compose -f "$compose" exec -T caddy wget -qO- http://localhost/healthz | grep -q '"status":"ok"'
# Assert restricted roles exist and cannot bypass row-level security.
docker compose -f "$compose" exec -T postgres sh -c 'PGPASSWORD="$POSTGRES_PASSWORD" psql -X -U "$POSTGRES_USER" -d "$POSTGRES_DB" -v ON_ERROR_STOP=1 -Atc "SELECT CASE WHEN count(*) = 2 THEN 1 ELSE 0 END FROM pg_roles WHERE rolname IN (\x27ploform_app\x27,\x27ploform_auth\x27) AND NOT rolbypassrls AND NOT rolsuper"' | grep -qx 1
# Confirm the API reports readiness (database PingContext).
docker compose -f "$compose" exec -T caddy wget -qO- http://api:8080/readyz >/dev/null
echo "Internal API and Caddy smoke checks passed"
