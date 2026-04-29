# Contributing to daraja-client-go

Thank you for your interest in contributing! We welcome contributions of all kinds.

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/YOUR-USERNAME/daraja-client-go`
3. Install Go 1.23 or later
4. Run `go mod download` to fetch dependencies
5. Create a feature branch: `git checkout -b feat/my-feature`

## Development Workflow

### Running Tests

```bash
go test ./... -race
go vet ./...
```

### Integration Tests

Integration tests hit the live sandbox and require credentials:

```bash
export DARAJA_CONSUMER_KEY=your_key
export DARAJA_CONSUMER_SECRET=your_secret
export DARAJA_SHORTCODE=174379
export DARAJA_PASSKEY=your_passkey
export DARAJA_INITIATOR=testapi
export DARAJA_SECURITY_CREDENTIAL=your_credential
go test -tags=integration ./...
```

### Linting

```bash
golangci-lint run ./...
```

### Security Scanning

```bash
go install golang.org/x/vuln/cmd/govulncheck@latest
govulncheck ./...
go install github.com/securego/gosec/v2/cmd/gosec@latest
gosec ./...
```

## Pull Request Process

1. Ensure all tests pass (`go test ./... -race`) and there are no lint/vet errors
2. Update documentation for any public API changes
3. Add a changelog entry for user-visible changes (see [CHANGELOG.md](CHANGELOG.md))
4. Fill out the [PR template](.github/pull_request_template.md) — it includes a security checklist
5. Request review from a maintainer

## Commit Conventions

Use [Conventional Commits](https://www.conventionalcommits.org/):

| Prefix | Usage |
|---|---|
| `feat:` | New feature |
| `fix:` | Bug fix |
| `docs:` | Documentation changes |
| `test:` | Test additions or changes |
| `refactor:` | Code refactoring |
| `chore:` | Tooling, CI, dependencies |
| `sec:` | Security improvements |

## Reporting Issues

- **Bug reports**: Use the [bug report template](.github/ISSUE_TEMPLATE/bug_report.yml)
- **Feature requests**: Use the [feature request template](.github/ISSUE_TEMPLATE/feature_request.yml)
- **Security vulnerabilities**: Follow [SECURITY.md](SECURITY.md) — do not open a public issue

## Code of Conduct

All contributors must follow our [Code of Conduct](CODE_OF_CONDUCT.md). Be respectful, inclusive, and constructive.

## Questions?

Open a [discussion](https://github.com/Kisotu/daraja-client-go/discussions) or check [SUPPORT.md](SUPPORT.md).