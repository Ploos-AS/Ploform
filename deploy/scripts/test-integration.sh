#!/bin/sh
# Run existing infrastructure checks, then transactional SQL integration tests.
set -eu
cd "$(dirname "$0")/../.."
sh deploy/scripts/smoke-test.sh
printf '\n==> PostgreSQL tenant isolation and form version integration tests\n'
docker compose -f deploy/compose.proxy.dev.yaml exec -T postgres \
  psql -X -U ploform_dev -d ploform -v ON_ERROR_STOP=1 \
  -f /dev/stdin < deploy/tests/tenant-isolation.sql

printf '\n==> PostgreSQL contact form persistence integration tests\n'
docker compose -f deploy/compose.proxy.dev.yaml exec -T postgres \
  psql -X -U ploform_dev -d ploform -v ON_ERROR_STOP=1 \
  -f /dev/stdin < deploy/tests/contact-form.sql

printf '\n==> Public HTTP submission end-to-end integration tests\n'
sh deploy/scripts/test-http-submission.sh
