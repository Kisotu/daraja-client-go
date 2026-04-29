// Package metrics provides observability metrics for the Daraja client.
package metrics

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestNewNoopRecorder(t *testing.T) {
	r := NewNoopRecorder()
	if r == nil {
		t.Fatal("NewNoopRecorder() returned nil")
	}
	// Should not panic
	ctx := context.Background()
	r.RecordRequestLatency(ctx, "/test", 100*time.Millisecond)
	r.IncrementRequestCounter(ctx, "/test", "200")
	r.IncrementErrorCounter(ctx, "/test", "network")
	r.IncrementTokenRefresh(true)
	r.SetCircuitBreakerState("payments", 0)
	r.IncInFlight("/test")
	r.DecInFlight("/test")
}

func TestRecordLatencyStartEnd(t *testing.T) {
	ctx := context.Background()
	ctx = RecordLatencyStart(ctx)

	// Small delay to ensure measurable time
	time.Sleep(10 * time.Millisecond)

	duration := RecordLatencyEnd(ctx)
	if duration < 10*time.Millisecond {
		t.Errorf("Expected at least 10ms, got %v", duration)
	}
}

func TestRecordLatencyEnd_NoStart(t *testing.T) {
	ctx := context.Background()
	duration := RecordLatencyEnd(ctx)
	if duration != 0 {
		t.Errorf("Expected 0 duration, got %v", duration)
	}
}

func TestSanitizeEndpoint(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"/api/v1/users", "api_v1_users"},
		{"/api/v1/users?id=123", "api_v1_users"},
		{"/", "root"},
		{"", "root"},
		{"/complex-path_name", "complex-path_name"},
		{"//api//test//", "api_test"},
	}

	for _, tt := range tests {
		result := SanitizeEndpoint(tt.input)
		if result != tt.expected {
			t.Errorf("SanitizeEndpoint(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestClassifyError(t *testing.T) {
	tests := []struct {
		err      error
		expected ErrorType
	}{
		{nil, ""},
		{errors.New("connection refused"), ErrorTypeNetwork},
		{errors.New("timeout occurred"), ErrorTypeTimeout},
		{errors.New("auth failed"), ErrorTypeAuth},
		{errors.New("circuit open"), ErrorTypeCircuitOpen},
		{errors.New("validation error"), ErrorTypeValidation},
		{errors.New("api returned 500"), ErrorTypeAPI},
		{errors.New("something else"), ErrorTypeUnknown},
	}

	for _, tt := range tests {
		result := ClassifyError(tt.err)
		if result != tt.expected {
			t.Errorf("ClassifyError(%v) = %v, want %v", tt.err, result, tt.expected)
		}
	}
}

func TestErrorType_String(t *testing.T) {
	if ErrorTypeNetwork.String() != "network" {
		t.Errorf("Expected 'network', got %q", ErrorTypeNetwork.String())
	}
}
