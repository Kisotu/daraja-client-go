// Package trace provides distributed tracing hooks for the Daraja client.
package trace

import (
	"context"
	"errors"
	"testing"
)

func TestNewNoopTracer(t *testing.T) {
	tracer := NewNoopTracer()
	if tracer == nil {
		t.Fatal("NewNoopTracer() returned nil")
	}

	ctx, span := tracer.Start(context.Background(), "test-span")
	if ctx == nil {
		t.Error("Start() returned nil context")
	}
	if span == nil {
		t.Error("Start() returned nil span")
	}

	// Should not panic
	span.End()
	span.SetStatus(StatusOK, "success")
	span.RecordError(errors.New("test error"))
	span.SetAttributes(String("key", "value"))
}

func TestNoopSpan(t *testing.T) {
	span := &noopSpan{}

	// All methods should be safe to call
	span.End()
	span.SetStatus(StatusOK, "ok")
	span.SetStatus(StatusError, "error")
	span.RecordError(errors.New("test"))
	span.SetAttributes(String("k", "v"), Int("i", 1))
}

func TestAttributeCreators(t *testing.T) {
	tests := []struct {
		name  string
		attr  Attribute
		key   string
		value interface{}
	}{
		{"String", String("key", "value"), "key", "value"},
		{"Int", Int("count", 42), "count", 42},
		{"Int64", Int64("size", 1000), "size", int64(1000)},
		{"Float64", Float64("ratio", 0.5), "ratio", 0.5},
		{"Bool", Bool("enabled", true), "enabled", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.attr.Key != tt.key {
				t.Errorf("Key = %q, want %q", tt.attr.Key, tt.key)
			}
			if tt.attr.Value != tt.value {
				t.Errorf("Value = %v, want %v", tt.attr.Value, tt.value)
			}
		})
	}
}

func TestSpanKind(t *testing.T) {
	kinds := []SpanKind{
		SpanKindInternal,
		SpanKindClient,
		SpanKindServer,
	}

	for _, kind := range kinds {
		// Just verify these don't panic
		_ = kind
	}
}

func TestStatus(t *testing.T) {
	statuses := []Status{
		StatusUnset,
		StatusOK,
		StatusError,
	}

	for _, status := range statuses {
		// Just verify these don't panic
		_ = status
	}
}

func TestContextWithSpan(t *testing.T) {
	ctx := context.Background()
	span := &noopSpan{}

	ctx = ContextWithSpan(ctx, span)
	retrieved := SpanFromContext(ctx)

	if retrieved != span {
		t.Error("Retrieved span should match the one stored")
	}
}

func TestSpanFromContext_NotFound(t *testing.T) {
	ctx := context.Background()
	span := SpanFromContext(ctx)

	if span != nil {
		t.Error("Expected nil span from empty context")
	}
}

func TestGlobalTracer(t *testing.T) {
	original := GlobalTracer

	// Set a new tracer
	newTracer := NewNoopTracer()
	SetGlobalTracer(newTracer)

	if GlobalTracer != newTracer {
		t.Error("Global tracer should have changed")
	}

	// Reset to nil (should become noop)
	SetGlobalTracer(nil)
	if _, ok := GlobalTracer.(*noopTracer); !ok {
		t.Error("Expected noop tracer after setting nil")
	}

	// Restore original
	SetGlobalTracer(original)
}

func TestSpanOptions(t *testing.T) {
	// Test that options don't panic
	var cfg spanConfig

	WithSpanKind(SpanKindClient)(&cfg)
	if cfg.kind != SpanKindClient {
		t.Error("Span kind not set correctly")
	}
}

func TestEventOptions(t *testing.T) {
	var cfg eventConfig

	WithStackTrace()(&cfg)
	if !cfg.stackTrace {
		t.Error("Stack trace not set")
	}

	attrs := []Attribute{String("key", "value")}
	WithEventAttributes(attrs...)(&cfg)
	if len(cfg.attrs) != 1 {
		t.Error("Attributes not set")
	}
}

func TestStartWithGlobal(t *testing.T) {
	ctx, span := Start(context.Background(), "test")
	if ctx == nil {
		t.Error("Start() returned nil context")
	}
	if span == nil {
		t.Error("Start() returned nil span")
	}
	span.End()
}

func BenchmarkNoopSpan(b *testing.B) {
	tracer := NewNoopTracer()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx, span := tracer.Start(ctx, "bench-span")
		span.SetAttributes(String("key", "value"))
		span.End()
		_ = ctx
	}
}
