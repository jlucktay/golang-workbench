#!/usr/bin/env bash
set -euo pipefail

echo "GET http://:8080/counter" \
  | vegeta attack -duration=30s \
  | tee "results.$(TZ=UTC date '+%Y%m%dT%H%M%SZ').bin" \
  | vegeta report
