// Package auth provides OAuth token management for the Daraja API.
package auth

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	defaultTokenBuffer = 60 * time.Second
)

// TokenResponse represents the OAuth token response from the Daraja API.
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   string `json:"expires_in"`
}

// AuthConfig holds the configuration for the AuthManager.
type AuthConfig struct {
	ConsumerKey    string
	ConsumerSecret string
	BaseURL        string
	HTTPClient     *http.Client
}

// AuthManager manages OAuth token retrieval, caching, and refresh.
type AuthManager struct {
	consumerKey    string
	consumerSecret string
	authURL        string
	httpClient     *http.Client

	mu          sync.RWMutex
	token       string
	expiresAt   time.Time
	refreshing  chan struct{}
	tokenBuffer time.Duration
}

// NewAuthManager creates a new AuthManager.
func NewAuthManager(cfg AuthConfig) *AuthManager {
	return &AuthManager{
		consumerKey:    cfg.ConsumerKey,
		consumerSecret: cfg.ConsumerSecret,
		authURL:        cfg.BaseURL + "/oauth/v1/generate?grant_type=client_credentials",
		httpClient:     cfg.HTTPClient,
		tokenBuffer:    defaultTokenBuffer,
	}
}

// Token returns a valid OAuth token, fetching a new one if necessary.
// It uses single-flight protection to prevent concurrent refreshes.
func (a *AuthManager) Token(ctx context.Context) (string, error) {
	a.mu.RLock()
	if a.token != "" && time.Now().Add(a.tokenBuffer).Before(a.expiresAt) {
		token := a.token
		a.mu.RUnlock()
		return token, nil
	}
	a.mu.RUnlock()

	return a.refreshToken(ctx)
}

// refreshToken fetches a new token with single-flight protection.
func (a *AuthManager) refreshToken(ctx context.Context) (string, error) {
	a.mu.Lock()

	if a.refreshing != nil {
		ch := a.refreshing
		a.mu.Unlock()
		select {
		case <-ch:
			a.mu.RLock()
			token := a.token
			a.mu.RUnlock()
			return token, nil
		case <-ctx.Done():
			return "", ctx.Err()
		}
	}

	a.refreshing = make(chan struct{})
	a.mu.Unlock()

	defer func() {
		a.mu.Lock()
		close(a.refreshing)
		a.refreshing = nil
		a.mu.Unlock()
	}()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, a.authURL, nil)
	if err != nil {
		return "", fmt.Errorf("auth: failed to create token request: %w", err)
	}

	auth := base64.StdEncoding.EncodeToString([]byte(a.consumerKey + ":" + a.consumerSecret))
	req.Header.Set("Authorization", "Basic "+auth)

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("auth: token request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("auth: failed to read token response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("auth: token request returned %d: %s", resp.StatusCode, body)
	}

	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", fmt.Errorf("auth: failed to unmarshal token response: %w", err)
	}

	token := tokenResp.AccessToken
	if token == "" {
		return "", fmt.Errorf("auth: empty access token in response")
	}

	expiresIn := 3599 * time.Second
	if tokenResp.ExpiresIn != "" {
		if d, err := time.ParseDuration(strings.TrimSpace(tokenResp.ExpiresIn) + "s"); err == nil {
			expiresIn = d
		}
	}

	a.mu.Lock()
	a.token = token
	a.expiresAt = time.Now().Add(expiresIn)
	a.mu.Unlock()

	return token, nil
}
