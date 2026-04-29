// Package trace provides distributed tracing hooks for the Daraja client.
// This package provides optional OpenTelemetry integration with safe no-op defaults.
package trace

import (
	"context"
)

// Status represents the status of a span.
type Status int

const (
	// StatusUnset is the default status.
	StatusUnset Status = iota
	// StatusOK indicates success.
	StatusOK
	// StatusError indicates an error.
	StatusError
)

// SpanKind represents the type of span.
type SpanKind int

const (
	// SpanKindInternal represents an internal operation.
	SpanKindInternal SpanKind = iota
	// SpanKindClient represents a client API call.
	SpanKindClient
	// SpanKindServer represents a server handler.
	SpanKindServer
)

// Attribute represents a span attribute.
type Attribute struct {
	Key   string
	Value interface{}
}

// String creates a string attribute.
func String(key, value string) Attribute {
	return Attribute{Key: key, Value: value}
}

// Int creates an int attribute.
func Int(key string, value int) Attribute {
	return Attribute{Key: key, Value: value}
}

// Int64 creates an int64 attribute.
func Int64(key string, value int64) Attribute {
	return Attribute{Key: key, Value: value}
}

// Float64 creates a float64 attribute.
func Float64(key string, value float64) Attribute {
	return Attribute{Key: key, Value: value}
}

// Bool creates a bool attribute.
func Bool(key string, value bool) Attribute {
	return Attribute{Key: key, Value: value}
}

// Span is the interface for a trace span.
type Span interface {
	// End completes the span.
	End()
	// SetStatus sets the span status.
	SetStatus(code Status, description string)
	// RecordError records an error on the span.
	RecordError(err error, options ...EventOption)
	// SetAttributes sets attributes on the span.
	SetAttributes(attrs ...Attribute)
}

// SpanOption configures span creation.
type SpanOption func(*spanConfig)

// EventOption configures event recording.
type EventOption func(*eventConfig)

type spanConfig struct {
	kind SpanKind
}

type eventConfig struct {
	stackTrace bool
	attrs      []Attribute
}

// WithSpanKind sets the span kind.
func WithSpanKind(kind SpanKind) SpanOption {
	return func(c *spanConfig) {
		c.kind = kind
	}
}

// WithStackTrace includes a stack trace with error recording.
func WithStackTrace() EventOption {
	return func(c *eventConfig) {
		c.stackTrace = true
	}
}

// WithEventAttributes sets attributes on the event.
func WithEventAttributes(attrs ...Attribute) EventOption {
	return func(c *eventConfig) {
		c.attrs = attrs
	}
}

// Tracer is the interface for starting spans.
type Tracer interface {
	// Start creates a new span and returns a context containing it.
	Start(ctx context.Context, spanName string, opts ...SpanOption) (context.Context, Span)
}

// contextKey is the key for storing spans in context.
type contextKey int

const spanKey contextKey = 0

// ContextWithSpan adds a span to the context.
func ContextWithSpan(ctx context.Context, span Span) context.Context {
	return context.WithValue(ctx, spanKey, span)
}

// SpanFromContext retrieves a span from the context.
// Returns nil if no span exists.
func SpanFromContext(ctx context.Context) Span {
	if span, ok := ctx.Value(spanKey).(Span); ok {
		return span
	}
	return nil
}

// noopSpan is a no-op implementation of Span.
type noopSpan struct{}

func (n *noopSpan) End()                                          {}
func (n *noopSpan) SetStatus(code Status, description string)     {}
func (n *noopSpan) RecordError(err error, options ...EventOption) {}
func (n *noopSpan) SetAttributes(attrs ...Attribute)              {}

// noopTracer is a no-op implementation of Tracer.
type noopTracer struct{}

func (n *noopTracer) Start(ctx context.Context, spanName string, opts ...SpanOption) (context.Context, Span) {
	return ctx, &noopSpan{}
}

// NewNoopTracer returns a no-op tracer that discards all spans.
// This is the default and is safe to use without initialization.
func NewNoopTracer() Tracer {
	return &noopTracer{}
}

// GlobalTracer is the package-level tracer.
var GlobalTracer Tracer = NewNoopTracer()

// SetGlobalTracer sets the global tracer.
func SetGlobalTracer(t Tracer) {
	if t == nil {
		GlobalTracer = NewNoopTracer()
	} else {
		GlobalTracer = t
	}
}

// Start starts a span using the global tracer.
func Start(ctx context.Context, spanName string, opts ...SpanOption) (context.Context, Span) {
	return GlobalTracer.Start(ctx, spanName, opts...)
}
