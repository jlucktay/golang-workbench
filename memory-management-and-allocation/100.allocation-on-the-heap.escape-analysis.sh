#!/usr/bin/env bash
set -euo pipefail

# Running the escape analysis command with 'go tool compile "-m" main.go' will confirm the allocation made by Go.
go tool compile "-m" main.go
