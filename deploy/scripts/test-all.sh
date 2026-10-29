#!/bin/sh
# Shared local/CI test entrypoint.
set -eu
cd "$(dirname "$0")/../.."

printf '\n==> Go unit tests\n'
(cd backend && go test ./...)

printf '\n==> Docker and PostgreSQL integration tests\n'
sh deploy/scripts/test-integration.sh

printf '\nAll configured checks passed.\n'
