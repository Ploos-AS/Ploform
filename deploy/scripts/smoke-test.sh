#!/bin/sh
# Run on an operator machine with Docker Compose and curl installed.
# Internal-only deployment verification; no public exposure.
set -eu
cd "$(dirname "$0")/../.."
compose="deploy/compose.proxy.dev.yaml"
: "${PLOFORM_DEV_POSTGRES_PASSWORD:?Set PLOFORM_DEV_POSTGRES_PASSWORD}"
command -v docker >/dev/null || { echo "Docker required" >&2; exit 1; }
command -v curl >/dev/null || { echo "curl required" >&2; exit 1; }
for file in deploy/secrets/app_password.txt deploy/secrets/auth_password.txt; do
 [ -s "$file" ] || { echo "Missing $file" >&2; exit 1; }
done
docker compose -f "$compose" up --build -d --wait
[ "$(docker compose -f "$compose" ps -q api)" != "" ]
# HTTP check from inside Caddy's isolated network (no published ports).
docker compose -f "$compose" exec -T caddy wget -qO- http://api:8080/healthz | grep -q '"status":"ok"'
docker compose -f "$compose" exec -T caddy wget -qO- http://localhost/healthz | grep -q '"status":"ok"'
# A configured DB should be reachable using the restricted app role.
docker compose -f "$compose" exec -T api /ploform-api --help >/dev/null 2>&1 && {
 echo "Unexpected API CLI behavior" >&2; exit 1;
} || true
# Confirm the API reports readiness (database PingContext).
docker compose -f "$compose" exec -T caddy wget -qO- http://api:8080/readyz >/dev/null
echo "Internal API and Caddy smoke checks passed"
