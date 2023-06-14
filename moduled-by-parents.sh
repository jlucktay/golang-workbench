#!/usr/bin/env bash
set -euo pipefail

if ! top_level=$(git rev-parse --show-toplevel); then
  echo >&2 "Not running from inside a git repo."
  exit 1
fi

# Find directories containing Go code that is not covered by a module that isn't the repo root.
fd --extension go --type file --exec bash -c 'echo {//}' "$top_level" \
  | sort -fu \
  | xargs -I % -n 1 bash -c 'if ! stat %/go.mod &> /dev/null; then ( echo %; cd %; go env GOMOD; ); fi' \
  | rg --before-context=1 --color=never --fixed-strings "$top_level"/go.mod
