// Package daraja is a production-ready Go client for Safaricom Daraja (M-Pesa) APIs.
package daraja

import (
	"context"
	"net/http"

	"github.com/anomalyco/daraja-client-go/auth"
	"github.com/anomalyco/daraja-client-go/transport"
)

// Client is the public interface for Daraja operations.
type Client interface {
	Token(ctx context.Context) (string, error)

	STKPush(ctx context.Context, req STKPushRequest) (*STKPushResponse, error)
}

type client struct {
	config  *Config
	authMgr *auth.AuthManager
	http    *http.Client
}

// NewClient creates a new Daraja client with the given options.
func NewClient(opts ...Option) (Client, error) {
	cfg := &Config{
		HTTPClient: transport.NewSecureClient(),
	}

	for _, opt := range opts {
		if err := opt(cfg); err != nil {
			return nil, err
		}
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	authMgr := auth.NewAuthManager(auth.AuthConfig{
		ConsumerKey:    cfg.ConsumerKey,
		ConsumerSecret: cfg.ConsumerSecret,
		BaseURL:        cfg.Environment.authURL(),
		HTTPClient:     cfg.HTTPClient,
	})

	return &client{
		config:  cfg,
		authMgr: authMgr,
		http:    cfg.HTTPClient,
	}, nil
}

func (c *client) Token(ctx context.Context) (string, error) {
	return c.authMgr.Token(ctx)
}