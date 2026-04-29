// Package logger provides structured logging with redaction of sensitive data.
package logger

import (
	"context"
	"log/slog"
	"os"
	"strings"
	"sync"
	"time"
)

// Logger is the logging interface used throughout the library.
// This abstraction allows for custom implementations while defaulting to no-op.
type Logger interface {
	// Debug logs a debug message with attributes
	Debug(msg string, attrs ...slog.Attr)
	// Info logs an info message with attributes
	Info(msg string, attrs ...slog.Attr)
	// Warn logs a warning message with attributes
	Warn(msg string, attrs ...slog.Attr)
	// Error logs an error message with attributes
	Error(msg string, attrs ...slog.Attr)
	// With returns a logger with the given attributes added
	With(attrs ...slog.Attr) Logger
}

// CorrelationIDKey is the context key type for correlation IDs.
type CorrelationIDKey struct{}

// RedactedString is the replacement value for sensitive data.
const RedactedString = "***REDACTED***"

// sensitiveFieldPatterns contains field names that trigger redaction.
// These are checked case-insensitively.
var sensitiveFieldPatterns = []string{
	"password",
	"passwd",
	"secret",
	"token",
	"apikey",
	"api_key",
	"authorization",
	"auth",
	"credential",
	"private_key",
	"passkey",
	"consumersecret",
	"consumer_secret",
}

// noopLogger is a no-op implementation of Logger.
type noopLogger struct{}

func (n *noopLogger) Debug(msg string, attrs ...slog.Attr) {}
func (n *noopLogger) Info(msg string, attrs ...slog.Attr)  {}
func (n *noopLogger) Warn(msg string, attrs ...slog.Attr)  {}
func (n *noopLogger) Error(msg string, attrs ...slog.Attr) {}
func (n *noopLogger) With(attrs ...slog.Attr) Logger       { return n }

// NewNoopLogger returns a no-op logger that discards all logs.
// This is the default and is safe to use without initialization.
func NewNoopLogger() Logger {
	return &noopLogger{}
}

// slogLogger wraps a standard library slog.Logger.
type slogLogger struct {
	logger *slog.Logger
	attrs  []slog.Attr
}

// NewSlogLogger wraps a slog.Logger for use as Logger.
func NewSlogLogger(l *slog.Logger) Logger {
	if l == nil {
		return NewNoopLogger()
	}
	return &slogLogger{logger: l}
}

// NewDefaultLogger creates a logger with sensible defaults for development.
func NewDefaultLogger(level slog.Level) Logger {
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: level,
	})
	return NewSlogLogger(slog.New(handler))
}

func (s *slogLogger) Debug(msg string, attrs ...slog.Attr) {
	attrs = s.withAttrs(attrs)
	attrs = RedactAttrs(attrs)
	s.logger.LogAttrs(context.TODO(), slog.LevelDebug, msg, attrs...)
}

func (s *slogLogger) Info(msg string, attrs ...slog.Attr) {
	attrs = s.withAttrs(attrs)
	attrs = RedactAttrs(attrs)
	s.logger.LogAttrs(context.TODO(), slog.LevelInfo, msg, attrs...)
}

func (s *slogLogger) Warn(msg string, attrs ...slog.Attr) {
	attrs = s.withAttrs(attrs)
	attrs = RedactAttrs(attrs)
	s.logger.LogAttrs(context.TODO(), slog.LevelWarn, msg, attrs...)
}

func (s *slogLogger) Error(msg string, attrs ...slog.Attr) {
	attrs = s.withAttrs(attrs)
	attrs = RedactAttrs(attrs)
	s.logger.LogAttrs(context.TODO(), slog.LevelError, msg, attrs...)
}

func (s *slogLogger) With(attrs ...slog.Attr) Logger {
	newAttrs := make([]slog.Attr, len(s.attrs)+len(attrs))
	copy(newAttrs, s.attrs)
	copy(newAttrs[len(s.attrs):], attrs)
	return &slogLogger{
		logger: s.logger,
		attrs:  newAttrs,
	}
}

func (s *slogLogger) withAttrs(attrs []slog.Attr) []slog.Attr {
	if len(s.attrs) == 0 {
		return attrs
	}
	result := make([]slog.Attr, len(s.attrs)+len(attrs))
	copy(result, s.attrs)
	copy(result[len(s.attrs):], attrs)
	return result
}

// RedactAttrs returns attrs with sensitive values redacted.
func RedactAttrs(attrs []slog.Attr) []slog.Attr {
	result := make([]slog.Attr, len(attrs))
	for i, attr := range attrs {
		result[i] = RedactAttr(attr)
	}
	return result
}

// RedactAttr redacts a single attribute if its key matches sensitive patterns.
func RedactAttr(attr slog.Attr) slog.Attr {
	if isSensitiveField(attr.Key) {
		return slog.String(attr.Key, RedactedString)
	}
	// Also check nested group attributes
	if attr.Value.Kind() == slog.KindGroup {
		groupAttrs := attr.Value.Group()
		redacted := RedactAttrs(groupAttrs)
		// Convert []slog.Attr to []any for slog.Group
		args := make([]any, len(redacted))
		for i, a := range redacted {
			args[i] = a
		}
		return slog.Group(attr.Key, args...)
	}
	return attr
}

// isSensitiveField checks if a field name indicates sensitive data.
func isSensitiveField(key string) bool {
	lowerKey := strings.ToLower(key)
	for _, pattern := range sensitiveFieldPatterns {
		if strings.Contains(lowerKey, pattern) {
			return true
		}
	}
	return false
}

// WithCorrelationID adds a correlation ID to the context.
func WithCorrelationID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, CorrelationIDKey{}, id)
}

// CorrelationIDFromContext extracts the correlation ID from context.
func CorrelationIDFromContext(ctx context.Context) string {
	if id, ok := ctx.Value(CorrelationIDKey{}).(string); ok {
		return id
	}
	return ""
}

// ContextWithCorrelationID returns attributes with correlation ID if present.
func ContextWithCorrelationID(ctx context.Context, attrs ...slog.Attr) []slog.Attr {
	if id := CorrelationIDFromContext(ctx); id != "" {
		return append([]slog.Attr{slog.String("correlation_id", id)}, attrs...)
	}
	return attrs
}

// RedactURL returns a URL string with credentials redacted.
func RedactURL(urlStr string) string {
	// Simple redaction of credentials in URLs
	// Removes user:pass@ from URLs
	if idx := strings.Index(urlStr, "@"); idx > 0 {
		// Find the scheme://
		if schemeIdx := strings.Index(urlStr, "://"); schemeIdx > 0 {
			scheme := urlStr[:schemeIdx+3]
			afterAt := urlStr[idx+1:]
			return scheme + "***CREDENTIALS***@" + afterAt
		}
	}
	return urlStr
}

// RedactMap returns a new map with sensitive values redacted.
func RedactMap(m map[string]string) map[string]string {
	if m == nil {
		return nil
	}
	result := make(map[string]string, len(m))
	for k, v := range m {
		if isSensitiveField(k) {
			result[k] = RedactedString
		} else {
			result[k] = v
		}
	}
	return result
}

// RedactJSON redacts sensitive values in a JSON string.
// Note: This is a simple implementation; for complex JSON use proper parsing.
func RedactJSON(jsonStr string) string {
	// This is a simple implementation that redacts common patterns
	result := jsonStr
	for _, pattern := range sensitiveFieldPatterns {
		// Pattern: "fieldName":"value"
		searchStr := `"` + pattern + `":`
		idx := 0
		for {
			idx = strings.Index(strings.ToLower(result[idx:]), searchStr)
			if idx == -1 {
				break
			}
			// Find the end of the value
			valueStart := idx + len(searchStr)
			// Skip whitespace
			for valueStart < len(result) && (result[valueStart] == ' ' || result[valueStart] == '\t') {
				valueStart++
			}
			if valueStart >= len(result) {
				break
			}

			var valueEnd int
			if result[valueStart] == '"' {
				// String value
				valueStart++ // Skip opening quote
				valueEnd = strings.Index(result[valueStart:], `"`)
				if valueEnd == -1 {
					break
				}
				valueEnd += valueStart
			} else {
				// Non-string value, find next comma or brace
				valueEnd = strings.IndexAny(result[valueStart:], `,}`)
				if valueEnd == -1 {
					break
				}
				valueEnd += valueStart
			}

			// Replace the value
			result = result[:valueStart] + RedactedString + result[valueEnd:]
			idx = valueStart + len(RedactedString)
		}
	}
	return result
}

// globalLogger is the package-level logger used by default.
var (
	globalLogger Logger = NewNoopLogger()
	loggerMu     sync.RWMutex
)

// SetGlobalLogger sets the global logger for the package.
func SetGlobalLogger(l Logger) {
	loggerMu.Lock()
	defer loggerMu.Unlock()
	if l == nil {
		globalLogger = NewNoopLogger()
	} else {
		globalLogger = l
	}
}

// Global returns the current global logger.
func Global() Logger {
	loggerMu.RLock()
	defer loggerMu.RUnlock()
	return globalLogger
}

// String creates a string attribute for logging.
func String(key, value string) slog.Attr {
	return slog.String(key, value)
}

// Int creates an int attribute for logging.
func Int(key string, value int) slog.Attr {
	return slog.Int(key, value)
}

// Duration creates a duration attribute for logging.
func Duration(key string, value time.Duration) slog.Attr {
	return slog.Duration(key, value)
}

// Any creates an attribute of any type.
func Any(key string, value interface{}) slog.Attr {
	return slog.Any(key, value)
}
