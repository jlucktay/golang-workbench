# Workspace Rules for go.jlucktay.dev/golang-workbench

## Module Creation and Structure
- **Module Path**: When initializing a new Go module, always use the module path `go.jlucktay.dev/golang-workbench/<full sub-path>`.
- **Nesting**: New modules should be nested appropriately under subdirectories based on their primary dependency or purpose.
- **go.work**: ALWAYS add newly created modules to the `go.work` file at the root of the workspace.

## Taskfile Operations
- Utilize the local `Taskfile.yaml` for formatting, linting, and other toil.
- When you create or modify code, use the existing Task targets:
  - `task fix`
  - `task modernise`
  - `task format`
  - `task lint`
  - `task toil` (runs all of the above)

## Preferred Libraries
- **Testing**: Use `github.com/matryer/is` for testing.
- **CLI Framework**: Use `github.com/urfave/cli` when building command-line tools.

## Coding Conventions
- **Table-Driven Tests**: Always stick to table-driven tests when writing test code.
- **Line Length (No Wrapping)**: DO NOT wrap long lines or columns in code, comments, prose, or anywhere, in any file format. No line length limit.
