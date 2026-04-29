//go:build otel

// Package trace provides OpenTelemetry integration.
// This file is only compiled when the otel build tag is present.
package trace

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// OTelTracer wraps an OpenTelemetry tracer.
type OTelTracer struct {
	tracer trace.Tracer
}

// NewOTelTracer creates a new OpenTelemetry-backed tracer.
func NewOTelTracer(tracer trace.Tracer) *OTelTracer {
	return &OTelTracer{tracer: tracer}
}

// Start creates a new span.
func (t *OTelTracer) Start(ctx context.Context, spanName string, opts ...SpanOption) (context.Context, Span) {
	var cfg spanConfig
	for _, opt := range opts {
		opt(&cfg)
	}

	// Convert SpanKind to OTel SpanKind
	var otelKind trace.SpanKind
	switch cfg.kind {
	case SpanKindClient:
		otelKind = trace.SpanKindClient
	case SpanKindServer:
		otelKind = trace.SpanKindServer
	default:
		otelKind = trace.SpanKindInternal
	}

	ctx, span := t.tracer.Start(ctx, spanName, trace.WithSpanKind(otelKind))
	return ctx, &otelSpan{span: span}
}

// otelSpan wraps an OpenTelemetry span.
type otelSpan struct {
	span trace.Span
}

// End completes the span.
func (s *otelSpan) End() {
	s.span.End()
}

// SetStatus sets the span status.
func (s *otelSpan) SetStatus(code Status, description string) {
	var otelCode codes.Code
	switch code {
	case StatusOK:
		otelCode = codes.Ok
	case StatusError:
		otelCode = codes.Error
	default:
		otelCode = codes.Unset
	}
	s.span.SetStatus(otelCode, description)
}

// RecordError records an error on the span.
func (s *otelSpan) RecordError(err error, options ...EventOption) {
	var cfg eventConfig
	for _, opt := range options {
		opt(&cfg)
	}

	otelOpts := make([]trace.EventOption, 0, len(cfg.attrs)+1)
	otelOpts = append(otelOpts, trace.WithStackTrace(cfg.stackTrace))

	for _, attr := range cfg.attrs {
		otelOpts = append(otelOpts, trace.WithAttributes(convertAttribute(attr)))
	}

	s.span.RecordError(err, otelOpts...)
}

// SetAttributes sets attributes on the span.
func (s *otelSpan) SetAttributes(attrs ...Attribute) {
	otelAttrs := make([]attribute.KeyValue, len(attrs))
	for i, attr := range attrs {
		otelAttrs[i] = convertAttribute(attr)
	}
	s.span.SetAttributes(otelAttrs...)
}

// convertAttribute converts our Attribute to OTel KeyValue.
func convertAttribute(attr Attribute) attribute.KeyValue {
	switch v := attr.Value.(type) {
	case string:
		return attribute.String(attr.Key, v)
	case int:
		return attribute.Int(attr.Key, v)
	case int64:
		return attribute.Int64(attr.Key, v)
	case float64:
		return attribute.Float64(attr.Key, v)
	case bool:
		return attribute.Bool(attr.Key, v)
	default:
		return attribute.String(attr.Key, "")
	}
}
