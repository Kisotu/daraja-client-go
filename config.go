package daraja

import (
	"errors"
	"net/http"
)

// Environment represents the Daraja API environment.
type Environment int

const (
	Sandbox Environment = iota
	Production
)

var (
	sandboxAuthURL    = "https://sandbox.safaricom.co.ke"
	productionAuthURL = "https://api.safaricom.co.ke"
)

func (e Environment) authURL() string {
	switch e {
	case Production:
		return productionAuthURL
	default:
		return sandboxAuthURL
	}
}

// Config holds the client configuration.
type Config struct {
	Environment    Environment
	ConsumerKey    string
	ConsumerSecret string
	ShortCode      string
	Passkey        string
	HTTPClient     *http.Client
}

// Option is a functional option for configuring the client.
type Option func(*Config) error

// WithEnvironment sets the Daraja API environment.
func WithEnvironment(env Environment) Option {
	return func(c *Config) error {
		c.Environment = env
		return nil
	}
}

// WithCredentials sets the Daraja API consumer key and secret.
func WithCredentials(consumerKey, consumerSecret string) Option {
	return func(c *Config) error {
		c.ConsumerKey = consumerKey
		c.ConsumerSecret = consumerSecret
		return nil
	}
}

// WithShortCode sets the business short code and passkey.
func WithShortCode(shortCode, passkey string) Option {
	return func(c *Config) error {
		c.ShortCode = shortCode
		c.Passkey = passkey
		return nil
	}
}

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Config) error {
		if httpClient == nil {
			return errors.New("daraja: HTTP client must not be nil")
		}
		c.HTTPClient = httpClient
		return nil
	}
}

func (c *Config) validate() error {
	if c.ConsumerKey == "" {
		return errors.New("daraja: consumer key is required")
	}
	if c.ConsumerSecret == "" {
		return errors.New("daraja: consumer secret is required")
	}
	return nil
}