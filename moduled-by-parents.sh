#!/usr/bin/env bash
set -euo pipefail

if ! top_level=$(git rev-parse --show-toplevel); then
  echo >&2 "Not running from inside a git repo."
  exit 1
fi

# Find directories containing Go code files which are included in a module, where that module is at the repo root.
fd --extension go --type file --base-directory "$top_level" --exec bash -c 'echo {//}' \
  | sort -fu \
  | sed "s|^\.|$top_level|" \
  | xargs -I % -S 1024 -n 1 bash -c 'if ! stat %/go.mod &> /dev/null; then ( echo %; go env -C % GOMOD; ); fi' \
  | rg --before-context=1 --color=never --fixed-strings "$top_level"/go.mod
