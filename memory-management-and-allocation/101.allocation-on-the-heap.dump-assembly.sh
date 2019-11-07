#!/usr/bin/env bash
set -euo pipefail

# Dumping the assembly code for this program, thanks to 'go tool compile -S main.go', would also explicitly show us the
# allocation.
go tool compile -S main.go | grep -A3 "smallStruct(SB), AX"

# The function 'newobject' is the built-in function for new allocations and proxy 'mallocgc', a function that manages
# them on the heap. There are two strategies in Go, one for the small allocations and one for larger ones.
