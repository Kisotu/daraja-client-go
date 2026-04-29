package auth

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func TestAuthManager_Token_Fetch(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/oauth/v1/generate" {
			http.NotFound(w, r)
			return
		}
		if r.Header.Get("Authorization") == "" {
			http.Error(w, "missing auth", http.StatusUnauthorized)
			return
		}
		resp := TokenResponse{
			AccessToken: "test-token",
			ExpiresIn:   "3599",
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	mgr := NewAuthManager(AuthConfig{
		ConsumerKey:    "key",
		ConsumerSecret: "secret",
		BaseURL:        srv.URL,
		HTTPClient:     srv.Client(),
	})

	token, err := mgr.Token(context.Background())
	if err != nil {
		t.Fatalf("Token() error = %v", err)
	}
	if token != "test-token" {
		t.Errorf("Token() = %q, want %q", token, "test-token")
	}
}

func TestAuthManager_Token_Cache(t *testing.T) {
	callCount := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		resp := TokenResponse{
			AccessToken: "test-token",
			ExpiresIn:   "3599",
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	mgr := NewAuthManager(AuthConfig{
		ConsumerKey:    "key",
		ConsumerSecret: "secret",
		BaseURL:        srv.URL,
		HTTPClient:     srv.Client(),
	})

	token1, err := mgr.Token(context.Background())
	if err != nil {
		t.Fatalf("Token() error = %v", err)
	}

	token2, err := mgr.Token(context.Background())
	if err != nil {
		t.Fatalf("Token() error = %v", err)
	}

	if token1 != token2 {
		t.Error("cached token should match first token")
	}
	if callCount != 1 {
		t.Errorf("expected 1 server call, got %d", callCount)
	}
}

func TestAuthManager_Token_ExpiredCache(t *testing.T) {
	callCount := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		resp := TokenResponse{
			AccessToken: "token-" + string(rune('A'+callCount-1)),
			ExpiresIn:   "1",
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	mgr := NewAuthManager(AuthConfig{
		ConsumerKey:    "key",
		ConsumerSecret: "secret",
		BaseURL:        srv.URL,
		HTTPClient:     srv.Client(),
	})
	mgr.tokenBuffer = 0

	_, err := mgr.Token(context.Background())
	if err != nil {
		t.Fatalf("Token() error = %v", err)
	}
	time.Sleep(2 * time.Second)

	_, err = mgr.Token(context.Background())
	if err != nil {
		t.Fatalf("Token() error = %v", err)
	}

	if callCount != 2 {
		t.Errorf("expected 2 server calls, got %d", callCount)
	}
}

func TestAuthManager_Token_SingleFlight(t *testing.T) {
	callCount := 0
	mu := sync.Mutex{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		callCount++
		mu.Unlock()
		time.Sleep(100 * time.Millisecond)
		resp := TokenResponse{
			AccessToken: "test-token",
			ExpiresIn:   "1",
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	mgr := NewAuthManager(AuthConfig{
		ConsumerKey:    "key",
		ConsumerSecret: "secret",
		BaseURL:        srv.URL,
		HTTPClient:     srv.Client(),
	})
	mgr.tokenBuffer = 0
	mgr.token = ""
	mgr.expiresAt = time.Time{}

	var wg sync.WaitGroup
	tokens := make([]string, 10)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			tok, err := mgr.Token(context.Background())
			if err != nil {
				t.Errorf("Token() error = %v", err)
				return
			}
			tokens[idx] = tok
		}(i)
	}
	wg.Wait()

	for i := 1; i < 10; i++ {
		if tokens[i] != tokens[0] {
			t.Errorf("all tokens should be equal, got %q and %q", tokens[0], tokens[i])
		}
	}
}

func TestAuthManager_Token_ServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "service unavailable", http.StatusServiceUnavailable)
	}))
	defer srv.Close()

	mgr := NewAuthManager(AuthConfig{
		ConsumerKey:    "key",
		ConsumerSecret: "secret",
		BaseURL:        srv.URL,
		HTTPClient:     srv.Client(),
	})

	_, err := mgr.Token(context.Background())
	if err == nil {
		t.Fatal("expected error from server error")
	}
}

func TestAuthManager_Token_EmptyToken(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := TokenResponse{
			AccessToken: "",
			ExpiresIn:   "3599",
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	mgr := NewAuthManager(AuthConfig{
		ConsumerKey:    "key",
		ConsumerSecret: "secret",
		BaseURL:        srv.URL,
		HTTPClient:     srv.Client(),
	})

	_, err := mgr.Token(context.Background())
	if err == nil {
		t.Fatal("expected error for empty token")
	}
}

func TestAuthManager_Token_ContextCancelled(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(500 * time.Millisecond)
	}))
	defer srv.Close()

	mgr := NewAuthManager(AuthConfig{
		ConsumerKey:    "key",
		ConsumerSecret: "secret",
		BaseURL:        srv.URL,
		HTTPClient:     srv.Client(),
	})
	mgr.tokenBuffer = 0
	mgr.token = ""
	mgr.expiresAt = time.Time{}

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	_, err := mgr.Token(ctx)
	if err == nil {
		t.Fatal("expected error from cancelled context")
	}
}

func TestAuthorizationHeader(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expected := "Basic " + base64.StdEncoding.EncodeToString([]byte("key:secret"))
		if r.Header.Get("Authorization") != expected {
			t.Errorf("Authorization = %q, want %s", r.Header.Get("Authorization"), expected)
		}
		resp := TokenResponse{
			AccessToken: "test-token",
			ExpiresIn:   "3599",
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	mgr := NewAuthManager(AuthConfig{
		ConsumerKey:    "key",
		ConsumerSecret: "secret",
		BaseURL:        srv.URL,
		HTTPClient:     srv.Client(),
	})

	_, err := mgr.Token(context.Background())
	if err != nil {
		t.Fatalf("Token() error = %v", err)
	}
}
