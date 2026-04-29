# ADR 0002: Retry with Jitter and Circuit Breaker

**Status:** Accepted  
**Date:** 2026-04-29  
**Author:** daraja-client-go team

## Context

The Daraja API can return transient errors (rate limiting, temporary outages).
The client must handle these gracefully without overwhelming the upstream service.

## Decision

Implement a two-layer resilience strategy:

1. **Retry middleware** (transport layer):
   - Exponential backoff with full jitter
   - Configurable max retries (default: 3)
   - Only retries on idempotent-safe conditions: 429, 5xx, network errors
   - Uses `crypto/rand` for jitter to avoid thundering herd

2. **Circuit breaker** (per-endpoint-group):
   - Three states: closed, open, half-open
   - Opens after configurable failure threshold
   - Half-open after configurable timeout
   - Closes after configurable success threshold
   - Metrics-exposed for monitoring

## Consequences

- + Upstream is protected from cascading failure
- + Transient errors are transparently handled
- + Circuit breaker prevents wasted calls during outages
- - Added complexity in debugging (transient retries mask some failures)