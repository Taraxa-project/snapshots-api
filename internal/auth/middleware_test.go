package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/taraxa/snapshots-api/internal/config"
)

func TestMiddleware_ExtractAPIKey(t *testing.T) {
	cfg := &config.Config{APIKeys: []string{"test-key"}}
	middleware := NewMiddleware(cfg)

	tests := []struct {
		name          string
		authHeader    string
		expectedKey   string
		expectedFound bool
	}{
		{
			name:          "valid bearer token",
			authHeader:    "Bearer test-key",
			expectedKey:   "test-key",
			expectedFound: true,
		},
		{
			name:          "valid bearer token lowercase",
			authHeader:    "bearer another-key",
			expectedKey:   "another-key",
			expectedFound: true,
		},
		{
			name:          "missing authorization header",
			authHeader:    "",
			expectedKey:   "",
			expectedFound: false,
		},
		{
			name:          "malformed header - no bearer",
			authHeader:    "test-key",
			expectedKey:   "",
			expectedFound: false,
		},
		{
			name:          "malformed header - only bearer",
			authHeader:    "Bearer",
			expectedKey:   "",
			expectedFound: false,
		},
		{
			name:          "different auth type",
			authHeader:    "Basic dGVzdDp0ZXN0",
			expectedKey:   "",
			expectedFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			key, found := middleware.ExtractAPIKey(req)

			if found != tt.expectedFound {
				t.Errorf("ExtractAPIKey() found = %v, want %v", found, tt.expectedFound)
			}

			if key != tt.expectedKey {
				t.Errorf("ExtractAPIKey() key = %v, want %v", key, tt.expectedKey)
			}
		})
	}
}

func TestMiddleware_IsAuthenticated(t *testing.T) {
	cfg := &config.Config{APIKeys: []string{"valid-key-1", "valid-key-2"}}
	middleware := NewMiddleware(cfg)

	tests := []struct {
		name           string
		authHeader     string
		expectedResult bool
	}{
		{
			name:           "valid API key 1",
			authHeader:     "Bearer valid-key-1",
			expectedResult: true,
		},
		{
			name:           "valid API key 2",
			authHeader:     "Bearer valid-key-2",
			expectedResult: true,
		},
		{
			name:           "invalid API key",
			authHeader:     "Bearer invalid-key",
			expectedResult: false,
		},
		{
			name:           "no authorization header",
			authHeader:     "",
			expectedResult: false,
		},
		{
			name:           "malformed header",
			authHeader:     "NotBearer valid-key-1",
			expectedResult: false,
		},
		{
			name:           "empty bearer token",
			authHeader:     "Bearer ",
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			result := middleware.IsAuthenticated(req)

			if result != tt.expectedResult {
				t.Errorf("IsAuthenticated() = %v, want %v", result, tt.expectedResult)
			}
		})
	}
}

func TestMiddleware_RequireAuth(t *testing.T) {
	cfg := &config.Config{APIKeys: []string{"valid-key"}}
	middleware := NewMiddleware(cfg)

	// Handler that should only be called for authenticated requests
	handlerCalled := false
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		w.WriteHeader(http.StatusOK)
	})

	wrappedHandler := middleware.RequireAuth(testHandler)

	tests := []struct {
		name                  string
		authHeader            string
		expectedStatus        int
		handlerShouldBeCalled bool
	}{
		{
			name:                  "valid API key - should call handler",
			authHeader:            "Bearer valid-key",
			expectedStatus:        http.StatusOK,
			handlerShouldBeCalled: true,
		},
		{
			name:                  "invalid API key - should reject",
			authHeader:            "Bearer invalid-key",
			expectedStatus:        http.StatusUnauthorized,
			handlerShouldBeCalled: false,
		},
		{
			name:                  "no auth header - should reject",
			authHeader:            "",
			expectedStatus:        http.StatusUnauthorized,
			handlerShouldBeCalled: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlerCalled = false // Reset

			req := httptest.NewRequest("GET", "/", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			rr := httptest.NewRecorder()
			wrappedHandler.ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("RequireAuth() status = %v, want %v", rr.Code, tt.expectedStatus)
			}

			if handlerCalled != tt.handlerShouldBeCalled {
				t.Errorf("RequireAuth() handler called = %v, want %v", handlerCalled, tt.handlerShouldBeCalled)
			}

			if tt.expectedStatus == http.StatusUnauthorized {
				// Check WWW-Authenticate header is set
				if auth := rr.Header().Get("WWW-Authenticate"); auth != "Bearer" {
					t.Errorf("Expected WWW-Authenticate header 'Bearer', got %v", auth)
				}

				// Check Content-Type is application/json
				if contentType := rr.Header().Get("Content-Type"); contentType != "application/json" {
					t.Errorf("Expected Content-Type 'application/json', got %v", contentType)
				}
			}
		})
	}
}

func TestConfig_IsValidAPIKey(t *testing.T) {
	cfg := &config.Config{
		APIKeys: []string{"key1", "key2", "key3"},
	}

	tests := []struct {
		name     string
		apiKey   string
		expected bool
	}{
		{
			name:     "valid key 1",
			apiKey:   "key1",
			expected: true,
		},
		{
			name:     "valid key 2",
			apiKey:   "key2",
			expected: true,
		},
		{
			name:     "valid key 3",
			apiKey:   "key3",
			expected: true,
		},
		{
			name:     "invalid key",
			apiKey:   "invalid-key",
			expected: false,
		},
		{
			name:     "empty key",
			apiKey:   "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cfg.IsValidAPIKey(tt.apiKey)
			if result != tt.expected {
				t.Errorf("IsValidAPIKey(%v) = %v, want %v", tt.apiKey, result, tt.expected)
			}
		})
	}
}

func TestConfig_IsValidAPIKey_EmptyConfig(t *testing.T) {
	cfg := &config.Config{
		APIKeys: []string{},
	}

	result := cfg.IsValidAPIKey("any-key")
	if result != false {
		t.Errorf("IsValidAPIKey() with empty config should return false, got %v", result)
	}
}
