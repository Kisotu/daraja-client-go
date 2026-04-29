# ADR 0001: Functional Options Pattern for Configuration

**Status:** Accepted  
**Date:** 2026-04-29  
**Author:** daraja-client-go team

## Context

The client needs a flexible configuration mechanism that:
- Supports optional settings with sensible defaults
- Remains backward-compatible as new options are added
- Is idiomatic Go

Alternatives considered: config struct, builder pattern, variadic options.

## Decision

Use the **functional options pattern**:

```go
type Option func(*Config) error

func WithTimeout(d time.Duration) Option {
    return func(c *Config) error {
        c.Timeout = d
        return nil
    }
}
```

This is the standard Go pattern popularized by Dave Cheney and used by
gRPC, AWS SDK Go v2, and other major Go libraries.

## Consequences

- + Backward compatible: new options don't change the function signature
- + Self-documenting: each option is a named function
- + Validation possible: options can return errors
- - More boilerplate than a simple config struct