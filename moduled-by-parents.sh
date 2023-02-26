#!/usr/bin/env bash
set -euo pipefail

# Find directories containing Go code that is not covered by a module that isn't the repo root.
fd --extension go --type file --exec bash -c 'echo {//}' \
  | sort -fu \
  | xargs -I % -n 1 bash -c 'if ! stat %/go.mod &> /dev/null; then ( echo %; cd %; go env GOMOD; ); fi' \
  | rg --before-context=1 --color=never --fixed-strings "$(pwd || true)/go.mod"
