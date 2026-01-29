# Contributing

Thanks for your interest in contributing to gocron. Please follow these guidelines to keep changes consistent and reviewable.

## Before you start
- Check existing issues and PRs to avoid duplication.
- For non-trivial changes, open an issue to discuss the approach first.

## Development setup
- Go version: see `go.mod`.
- Install dependencies with standard Go tooling.

## Coding standards
- Keep changes minimal and focused.
- Follow existing naming and structure conventions.
- Prefer explicit error handling; avoid hidden side effects.

## Testing and linting
- Unit testing guidance: see `ai-rules/test/TESTS.md`.
- After code changes, run:
  - `goimports -w -l ./...`
  - `golangci-lint run`

## Pull requests
- Keep PRs small and focused.
- Describe the change, motivation, and any trade-offs.
- Link relevant issues and include test results when applicable.

## Reporting issues
- Provide reproduction steps and expected vs actual behavior.
- Include environment details (OS, Go version).
