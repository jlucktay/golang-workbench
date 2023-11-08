#!/usr/bin/env bash
set -euo pipefail

# If there was already a 'go.mod' file here, this will be an update and not a net-new setup.
if [[ -f go.mod ]]; then
  net_new_mod=0
  commit_message_suffix="update module"
else
  net_new_mod=1
  commit_message_suffix="set up a module"
fi

if [[ -n ${1-} ]]; then
  # If an argument is given, clean slate.
  rm -fv go.{mod,sum}
fi

# Get the current subdirectory within the git repo, sans trailing slash.
this_module=$(git rev-parse --show-prefix)
this_module=${this_module%/}

if [[ -n $this_module ]]; then
  go work edit -use="./$this_module"
  this_module="($this_module)"
fi

if [[ $net_new_mod == 1 ]]; then
  init_module=$(pwd | cut -d'/' -f5-)
  go mod init "$init_module"
fi

go mod tidy -go=1.21
go vet ./...
go build ./...
go test --count=1 -v ./...
go clean -x

git add ./go.*

if [[ $net_new_mod == 1 ]]; then
  toplevel=$(git rev-parse --show-toplevel)
  git ap -- "$toplevel"/go.work
fi

git commit --message="build$this_module: $commit_message_suffix"
