# Contributing to daraja-client-go

Thank you for your interest in contributing!

## Development Setup

1. Install Go 1.23 or later
2. Clone the repository
3. Run `go mod download`

## Running Tests

```bash
go test ./... -race
go vet ./...
```

## Pull Request Process

1. Ensure tests pass and no lint/vet errors
2. Update documentation for any public API changes
3. Add a changelog entry for user-visible changes

## Commit Conventions

Use [Conventional Commits](https://www.conventionalcommits.org/) format:
- `feat:` for new features
- `fix:` for bug fixes
- `docs:` for documentation changes
- `test:` for test additions/changes

## Code of Conduct

See [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md).