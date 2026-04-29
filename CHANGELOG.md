# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- CLI scaffold (`cmd/daraja-cli`) with auth, STK push, and status query commands
- Complete documentation set: getting started, configuration, security,
  error handling, examples, release process, go-live checklist
- Architecture Decision Records (ADRs) for functional options, retry/circuit
  breaker, and observability interfaces
- Mermaid sequence diagrams for OAuth lifecycle, STK Push callback, and
  Reversal workflow
- Pull request template and issue templates (bug report, feature request,
  security concern)
- Expanded CONTRIBUTING.md with full development workflow and conventions
- README badges (CI, Go Reference, Go Report Card, License)
- Concurrency and stress tests for STK Push, C2B, B2C, Transaction Status,
  and Reversal under high contention
- Performance benchmarks for STK Push, STK Push Query, C2B Simulate,
  token refresh, JSON marshal/unmarshal, and parallel execution
- API compatibility checklist documenting the frozen public surface
- Release process documentation and GitHub release automation workflow
- CHANGELOG.md following Keep a Changelog format

### Changed

- README rewritten with comprehensive feature list, CLI section, API coverage
  table, and links to all documentation

## [0.9.0] - 2026-04-27

### Added

- STK Push, C2B, B2C, Transaction Status, Reversal API endpoints
- OAuth token management with caching and single-flight refresh
- Secure HTTP client with TLS 1.2+, timeouts, and connection pooling
- Retry middleware with exponential backoff and jitter
- Circuit breaker pattern per endpoint group
- Structured logging with secret redaction
- Prometheus metrics counters and histograms
- Optional OpenTelemetry tracing integration
- Callback verification utilities (IP allowlist, payload validation,
  correlation ID matching)
- Sandbox and production environment support
- Typed DarajaError model with retryability detection
- Functional options pattern for client configuration
- CI pipeline with lint, test, security scanning
- Integration tests behind build tag
- Open-source repository files (LICENSE, CONTRIBUTING, CODE_OF_CONDUCT,
  SECURITY, SUPPORT)

[Unreleased]: https://github.com/Kisotu/daraja-client-go/compare/v0.9.0...HEAD
[0.9.0]: https://github.com/Kisotu/daraja-client-go/releases/tag/v0.9.0