#!/usr/bin/env bash
set -euo pipefail

while true; do
  curl --location http://localhost:8080
  echo
done
