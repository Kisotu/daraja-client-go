package transport

import (
	"crypto/tls"
	"net/http"
	"testing"
)

func TestNewSecureClient(t *testing.T) {
	client := NewSecureClient()
	if client == nil {
		t.Fatal("NewSecureClient() = nil")
	}
	if client.Timeout == 0 {
		t.Error("client has no timeout")
	}

	transport, ok := client.Transport.(*http.Transport)
	if !ok {
		t.Fatal("Transport is not *http.Transport")
	}
	if transport.TLSClientConfig == nil {
		t.Fatal("TLSClientConfig is nil")
	}
	if transport.TLSClientConfig.MinVersion != tls.VersionTLS12 {
		t.Errorf("MinVersion = %v, want TLS 1.2", transport.TLSClientConfig.MinVersion)
	}
	if transport.MaxIdleConns == 0 {
		t.Error("MaxIdleConns is zero")
	}
	if transport.MaxIdleConnsPerHost == 0 {
		t.Error("MaxIdleConnsPerHost is zero")
	}
	if transport.IdleConnTimeout == 0 {
		t.Error("IdleConnTimeout is zero")
	}
}