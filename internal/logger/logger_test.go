// Package logger provides structured logging with redaction of sensitive data.
package logger

import (
	"context"
	"log/slog"
	"strings"
	"testing"
)

func TestNewNoopLogger(t *testing.T) {
	l := NewNoopLogger()
	if l == nil {
		t.Fatal("NewNoopLogger() returned nil")
	}
	// Should not panic
	l.Debug("debug message")
	l.Info("info message")
	l.Warn("warn message")
	l.Error("error message")
	l2 := l.With(slog.String("key", "value"))
	if l2 == nil {
		t.Error("With() returned nil")
	}
}

func TestNewSlogLogger(t *testing.T) {
	l := slog.New(slog.DiscardHandler)
	logger := NewSlogLogger(l)
	if logger == nil {
		t.Fatal("NewSlogLogger() returned nil")
	}
}

func TestNewSlogLogger_Nil(t *testing.T) {
	logger := NewSlogLogger(nil)
	if logger == nil {
		t.Fatal("Expected non-nil logger")
	}
	// Should be noop
	logger.Info("test")
}

func TestNewDefaultLogger(t *testing.T) {
	logger := NewDefaultLogger(slog.LevelDebug)
	if logger == nil {
		t.Fatal("NewDefaultLogger() returned nil")
	}
}

func TestRedactAttr_Sensitive(t *testing.T) {
	tests := []struct {
		key      string
		value    string
		expected string
	}{
		{"password", "secret123", RedactedString},
		{"Password", "secret123", RedactedString},
		{"api_key", "abc123", RedactedString},
		{"Authorization", "Bearer token", RedactedString},
		{"passkey", "mypasskey", RedactedString},
		{"consumer_secret", "secret", RedactedString},
		{"normal_field", "value", "value"},
		{"user_name", "john", "john"},
	}

	for _, tt := range tests {
		attr := slog.String(tt.key, tt.value)
		result := RedactAttr(attr)
		if result.Value.String() != tt.expected {
			t.Errorf("RedactAttr(%q) = %q, want %q", tt.key, result.Value.String(), tt.expected)
		}
	}
}

func TestRedactAttrs(t *testing.T) {
	attrs := []slog.Attr{
		slog.String("user", "john"),
		slog.String("password", "secret"),
		slog.String("age", "30"),
		slog.String("api_key", "key123"),
	}

	result := RedactAttrs(attrs)

	if result[0].Value.String() != "john" {
		t.Errorf("user field should not be redacted, got %q", result[0].Value.String())
	}
	if result[1].Value.String() != RedactedString {
		t.Errorf("password field should be redacted, got %q", result[1].Value.String())
	}
	if result[2].Value.String() != "30" {
		t.Errorf("age field should not be redacted, got %q", result[2].Value.String())
	}
	if result[3].Value.String() != RedactedString {
		t.Errorf("api_key field should be redacted, got %q", result[3].Value.String())
	}
}

func TestIsSensitiveField(t *testing.T) {
	sensitive := []string{"password", "PASSWORD", "Password", "api_key", "Authorization", "token", "consumerSecret"}
	notSensitive := []string{"user", "name", "email", "timestamp"}

	for _, field := range sensitive {
		if !isSensitiveField(field) {
			t.Errorf("isSensitiveField(%q) should be true", field)
		}
	}

	for _, field := range notSensitive {
		if isSensitiveField(field) {
			t.Errorf("isSensitiveField(%q) should be false", field)
		}
	}
}

func TestWithCorrelationID(t *testing.T) {
	ctx := context.Background()
	ctx = WithCorrelationID(ctx, "abc-123")
	id := CorrelationIDFromContext(ctx)
	if id != "abc-123" {
		t.Errorf("CorrelationIDFromContext() = %q, want %q", id, "abc-123")
	}
}

func TestCorrelationIDFromContext_Missing(t *testing.T) {
	ctx := context.Background()
	id := CorrelationIDFromContext(ctx)
	if id != "" {
		t.Errorf("Expected empty string, got %q", id)
	}
}

func TestContextWithCorrelationID(t *testing.T) {
	ctx := WithCorrelationID(context.Background(), "corr-123")
	attrs := ContextWithCorrelationID(ctx, slog.String("key", "value"))

	if len(attrs) != 2 {
		t.Fatalf("Expected 2 attributes, got %d", len(attrs))
	}
	if attrs[0].Key != "correlation_id" || attrs[0].Value.String() != "corr-123" {
		t.Errorf("First attr should be correlation_id, got %v", attrs[0])
	}
}

func TestContextWithCorrelationID_NoID(t *testing.T) {
	ctx := context.Background()
	attrs := ContextWithCorrelationID(ctx, slog.String("key", "value"))

	if len(attrs) != 1 {
		t.Fatalf("Expected 1 attribute, got %d", len(attrs))
	}
	if attrs[0].Key != "key" {
		t.Errorf("Expected key attr, got %v", attrs[0])
	}
}

func TestRedactURL(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"https://user:pass@example.com/path", "https://***CREDENTIALS***@example.com/path"},
		{"http://admin:secret@api.com", "http://***CREDENTIALS***@api.com"},
		{"https://example.com/path", "https://example.com/path"},
		{"http://localhost:8080", "http://localhost:8080"},
	}

	for _, tt := range tests {
		result := RedactURL(tt.input)
		if result != tt.expected {
			t.Errorf("RedactURL(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestRedactMap(t *testing.T) {
	input := map[string]string{
		"user":         "john",
		"password":     "secret",
		"api_key":      "key123",
		"normal_field": "value",
	}

	result := RedactMap(input)

	if result["user"] != "john" {
		t.Errorf("user should not be redacted")
	}
	if result["password"] != RedactedString {
		t.Errorf("password should be redacted, got %q", result["password"])
	}
	if result["api_key"] != RedactedString {
		t.Errorf("api_key should be redacted, got %q", result["api_key"])
	}
	if result["normal_field"] != "value" {
		t.Errorf("normal_field should not be redacted")
	}
}

func TestRedactMap_Nil(t *testing.T) {
	result := RedactMap(nil)
	if result != nil {
		t.Errorf("Expected nil, got %v", result)
	}
}

func TestRedactJSON(t *testing.T) {
	input := `{"user": "john", "password": "secret", "api_key": "xyz123", "age": 30}`
	result := RedactJSON(input)

	if !strings.Contains(result, `"user":`) {
		t.Error("user field should remain")
	}
	if !strings.Contains(result, RedactedString) {
		t.Errorf("sensitive values should be redacted, got: %s", result)
	}
}

func TestSlogLogger_With(t *testing.T) {
	base := slog.New(slog.DiscardHandler)
	logger := NewSlogLogger(base).(*slogLogger)

	withAttrs := logger.With(slog.String("key1", "value1"))
	if withAttrs == nil {
		t.Fatal("With() returned nil")
	}

	// Should return a new logger
	if withAttrs == logger {
		t.Error("With() should return a new logger instance")
	}
}

func TestGlobalLogger(t *testing.T) {
	original := Global()

	// Set a new logger
	newLogger := NewDefaultLogger(slog.LevelInfo)
	SetGlobalLogger(newLogger)

	if Global() == original {
		t.Error("Global logger should have changed")
	}

	// Reset to noop
	SetGlobalLogger(nil)
	if _, ok := Global().(*noopLogger); !ok {
		t.Error("Expected noop logger after setting nil")
	}

	// Restore original
	SetGlobalLogger(original)
}

func TestHelperFunctions(t *testing.T) {
	// Test String helper
	attr := String("key", "value")
	if attr.Key != "key" || attr.Value.String() != "value" {
		t.Error("String helper failed")
	}

	// Test Int helper
	attr = Int("count", 42)
	if attr.Key != "count" || attr.Value.Int64() != 42 {
		t.Error("Int helper failed")
	}

	// Test Duration helper - we just verify it doesn't panic
	_ = Duration("elapsed", 0)

	// Test Any helper
	attr = Any("data", map[string]string{"key": "value"})
	if attr.Key != "data" {
		t.Error("Any helper failed")
	}
}
