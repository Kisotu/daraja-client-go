package daraja

import (
	"errors"
	"net/http"
	"testing"
)

func TestNewClient_ValidConfig(t *testing.T) {
	client, err := NewClient(
		WithEnvironment(Sandbox),
		WithCredentials("key", "secret"),
	)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	if client == nil {
		t.Fatal("NewClient() = nil")
	}
}

func TestNewClient_MissingConsumerKey(t *testing.T) {
	_, err := NewClient(
		WithCredentials("", "secret"),
	)
	if err == nil {
		t.Fatal("expected error for missing consumer key")
	}
}

func TestNewClient_MissingConsumerSecret(t *testing.T) {
	_, err := NewClient(
		WithCredentials("key", ""),
	)
	if err == nil {
		t.Fatal("expected error for missing consumer secret")
	}
}

func TestNewClient_NilHTTPClient(t *testing.T) {
	_, err := NewClient(
		WithCredentials("key", "secret"),
		WithHTTPClient(nil),
	)
	if err == nil {
		t.Fatal("expected error for nil HTTP client")
	}
}

func TestNewClient_CustomHTTPClient(t *testing.T) {
	custom := &http.Client{}
	c, err := NewClient(
		WithCredentials("key", "secret"),
		WithHTTPClient(custom),
	)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	if c.(*client).http != custom {
		t.Fatal("custom HTTP client not used")
	}
}

func TestEnvironment_AuthURL(t *testing.T) {
	tests := []struct {
		name string
		env  Environment
		want string
	}{
		{"sandbox", Sandbox, "https://sandbox.safaricom.co.ke"},
		{"production", Production, "https://api.safaricom.co.ke"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.env.authURL(); got != tt.want {
				t.Errorf("authURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDarajaError(t *testing.T) {
	rootErr := errors.New("connection refused")
	de := &DarajaError{
		HTTPStatus:   500,
		ResponseCode: "500.001.1001",
		ResponseDesc: "Internal server error",
		Endpoint:     "/mpesa/stkpush/v1/processrequest",
		Retryable:    true,
		Err:          rootErr,
	}

	errMsg := de.Error()
	if errMsg == "" {
		t.Fatal("Error() returned empty string")
	}

	if !errors.Is(de, rootErr) {
		t.Fatal("Unwrap() did not return root error")
	}
}
