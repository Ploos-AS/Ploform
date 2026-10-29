#!/bin/sh
# Local test entrypoint; GitHub Actions can invoke the same script.
set -eu
cd "$(dirname "$0")/../.."

printf '\n==> Go unit tests\n'
(cd backend && go test ./...)

printf '\n==> Docker/PostgreSQL smoke tests\n'
sh deploy/scripts/smoke-test.sh

printf '\nAll currently implemented checks passed.\n'
printf 'NOTE: Cross-tenant RLS and form-version integration tests are not yet implemented.\n'
