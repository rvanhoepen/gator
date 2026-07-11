#!/usr/bin/env bash
set -euo pipefail

if [ -f .env.test ]; then
  set -a
  source .env.test
  set +a
fi

: "${TEST_DATABASE_URL:?TEST_DATABASE_URL is required}"

if [ "${1:-}" = "coveragehtml" ] || [ "${1:-}" = "coverage" ]; then
  go test -v -coverprofile=coverage.out ./...
  go tool cover -html=coverage.out -o coverage.html
  firefox coverage.html
else
  go test -v ./...
fi
