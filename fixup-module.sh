#!/usr/bin/env bash
set -euo pipefail

# If there was already a 'go.mod' file here, this will be an update and not a net-new setup.
if [[ -f go.mod ]]; then
  commit_message_suffix="update a module"
else
  commit_message_suffix="set up a module"
fi

# Either way, clean slate.
rm -fv go.{mod,sum}

# Get the current subdirectory within the git repo, sans trailing slash.
this_module=$(git rev-parse --show-prefix)
this_module=${this_module%/}

if [[ -n $this_module ]]; then
  go work edit -use="./$this_module"
  this_module="($this_module)"
fi

init_module=$(pwd | cut -d'/' -f5-)
go mod init "$init_module"

go mod tidy
go vet ./...
go build ./...
go test --count=1 -v ./...
go clean -x

git add ./go.*
git ap
git commit --message="build$this_module: $commit_message_suffix"
