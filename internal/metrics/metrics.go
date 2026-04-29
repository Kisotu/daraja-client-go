// Package metrics provides observability metrics for the Daraja client.
package metrics

import (
	"context"
	"strings"
	"time"
)

// Recorder is the metrics interface used throughout the library.
// Implementations can be backed by Prometheus, StatsD, or be no-op.
type Recorder interface {
	// RecordRequestLatency records the duration of a request
	RecordRequestLatency(ctx context.Context, endpoint string, duration time.Duration)
	// IncrementRequestCounter increments the total request count
	IncrementRequestCounter(ctx context.Context, endpoint, status string)
	// IncrementErrorCounter increments the error count by type
	IncrementErrorCounter(ctx context.Context, endpoint, errorType string)
	// IncrementTokenRefresh records token refresh events
	IncrementTokenRefresh(success bool)
	// SetCircuitBreakerState updates the circuit breaker state gauge
	SetCircuitBreakerState(endpointGroup string, state int)
	// IncInFlight increments in-flight request count
	IncInFlight(endpoint string)
	// DecInFlight decrements in-flight request count
	DecInFlight(endpoint string)
}

// noopRecorder is a no-op implementation of Recorder.
type noopRecorder struct{}

func (n *noopRecorder) RecordRequestLatency(ctx context.Context, endpoint string, duration time.Duration) {
}
func (n *noopRecorder) IncrementRequestCounter(ctx context.Context, endpoint, status string)  {}
func (n *noopRecorder) IncrementErrorCounter(ctx context.Context, endpoint, errorType string) {}
func (n *noopRecorder) IncrementTokenRefresh(success bool)                                    {}
func (n *noopRecorder) SetCircuitBreakerState(endpointGroup string, state int)                {}
func (n *noopRecorder) IncInFlight(endpoint string)                                           {}
func (n *noopRecorder) DecInFlight(endpoint string)                                           {}

// NewNoopRecorder returns a no-op Recorder that discards all metrics.
func NewNoopRecorder() Recorder {
	return &noopRecorder{}
}

// contextKey is used for storing metrics-related values in context.
type contextKey int

const (
	startTimeKey contextKey = iota
)

// RecordLatencyStart returns a context with the start time recorded.
// Use with RecordLatencyEnd to measure request duration.
func RecordLatencyStart(ctx context.Context) context.Context {
	return context.WithValue(ctx, startTimeKey, time.Now())
}

// RecordLatencyEnd calculates and returns the elapsed duration since start.
// Returns zero duration if start time was not recorded.
func RecordLatencyEnd(ctx context.Context) time.Duration {
	if start, ok := ctx.Value(startTimeKey).(time.Time); ok {
		return time.Since(start)
	}
	return 0
}

// SanitizeEndpoint creates a safe metric label from an endpoint path.
// Removes query parameters, replaces special characters.
func SanitizeEndpoint(endpoint string) string {
	// Remove query parameters
	for i := 0; i < len(endpoint); i++ {
		if endpoint[i] == '?' {
			endpoint = endpoint[:i]
			break
		}
	}

	// Replace special characters with underscores
	result := make([]byte, 0, len(endpoint))
	for i := 0; i < len(endpoint); i++ {
		c := endpoint[i]
		switch {
		case (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_' || c == '-':
			result = append(result, c)
		case c == '/':
			if len(result) > 0 && result[len(result)-1] != '_' {
				result = append(result, '_')
			}
		default:
			if len(result) > 0 && result[len(result)-1] != '_' {
				result = append(result, '_')
			}
		}
	}

	// Trim leading/trailing underscores
	s := string(result)
	for len(s) > 0 && s[0] == '_' {
		s = s[1:]
	}
	for len(s) > 0 && s[len(s)-1] == '_' {
		s = s[:len(s)-1]
	}

	if s == "" {
		return "root"
	}
	return s
}

// ErrorType categorizes errors for metrics.
type ErrorType string

const (
	ErrorTypeNetwork     ErrorType = "network"
	ErrorTypeTimeout     ErrorType = "timeout"
	ErrorTypeAuth        ErrorType = "auth"
	ErrorTypeValidation  ErrorType = "validation"
	ErrorTypeAPI         ErrorType = "api"
	ErrorTypeCircuitOpen ErrorType = "circuit_open"
	ErrorTypeUnknown     ErrorType = "unknown"
)

// String returns the string representation.
func (e ErrorType) String() string {
	return string(e)
}

// ClassifyError determines the error type for metrics.
func ClassifyError(err error) ErrorType {
	if err == nil {
		return ""
	}

	errStr := strings.ToLower(err.Error())

	// Check for specific error patterns
	switch {
	case strings.Contains(errStr, "circuit") && strings.Contains(errStr, "open"):
		return ErrorTypeCircuitOpen
	case strings.Contains(errStr, "timeout") || strings.Contains(errStr, "deadline"):
		return ErrorTypeTimeout
	case strings.Contains(errStr, "auth") || strings.Contains(errStr, "token") || strings.Contains(errStr, "unauthorized"):
		return ErrorTypeAuth
	case strings.Contains(errStr, "connection") || strings.Contains(errStr, "network") || strings.Contains(errStr, "dial"):
		return ErrorTypeNetwork
	case strings.Contains(errStr, "validation") || strings.Contains(errStr, "invalid"):
		return ErrorTypeValidation
	case strings.Contains(errStr, "api") || strings.Contains(errStr, "response"):
		return ErrorTypeAPI
	default:
		return ErrorTypeUnknown
	}
}
