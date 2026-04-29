# ADR 0003: Observability via Interfaces with No-Op Defaults

**Status:** Accepted  
**Date:** 2026-04-29  
**Author:** daraja-client-go team

## Context

The client must support logging, metrics, and tracing without imposing
specific implementations on consumers. Some users will want Prometheus,
others will use Datadog or another provider. The client must not force
a dependency on any specific observability backend.

## Decision

Define **interfaces** in internal packages with **no-op default implementations**:

- `internal/logger/Logger` — structured logging with severity levels
- `internal/metrics/Recorder` — counters, histograms, gauges
- `internal/trace/Tracer` — span creation and propagation

Users supply implementations via options (e.g. `WithLogger()`). If none
is provided, the no-op implementation silently discards all events.

For tracing, OpenTelemetry is supported via a build tag (`otel`) to avoid
unnecessary dependency bloat for users who don't need it.

## Consequences

- + Zero dependency overhead for users who don't need observability
- + Freedom to use any backend (Prometheus, Datadog, OpenTelemetry, etc.)
- + Testable: no-op implementations simplify unit testing
- - Slightly more internal plumbing compared to hard-coded logging