#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")/.."

if ! command -v go >/dev/null 2>&1; then
  echo "Go 1.23+ is required" >&2
  exit 1
fi

echo "Checking Go formatting..."
unformatted="$(gofmt -l backend)"
if [[ -n "$unformatted" ]]; then
  printf 'Unformatted Go files:\n%s\n' "$unformatted" >&2
  exit 1
fi

echo "Running Go tests (without changing go.mod/go.sum)..."
(cd backend && go test -mod=readonly ./...)

echo "Running go vet..."
(cd backend && go vet -mod=readonly ./...)

echo "Go checks passed."
