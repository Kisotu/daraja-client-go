// Package transport provides secure HTTP client configuration and middleware.
package transport

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"
)

// NewSecureClient creates an *http.Client with production-safe defaults.
func NewSecureClient() *http.Client {
	t := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   10,
		IdleConnTimeout:       90 * time.Second,
		TLSClientConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	}

	return &http.Client{
		Transport: t,
		Timeout:   30 * time.Second,
	}
}