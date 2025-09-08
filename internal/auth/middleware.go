package auth

import (
	"net/http"
	"strings"

	"github.com/taraxa/snapshots-api/internal/config"
)

// Middleware provides authentication functionality
type Middleware struct {
	config *config.Config
}

// NewMiddleware creates a new authentication middleware
func NewMiddleware(cfg *config.Config) *Middleware {
	return &Middleware{
		config: cfg,
	}
}

// ExtractAPIKey extracts API key from Authorization header
// Returns the API key and whether it was found
func (m *Middleware) ExtractAPIKey(r *http.Request) (string, bool) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", false
	}

	// Check for Bearer token format: "Bearer <token>"
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", false
	}

	return parts[1], true
}

// IsAuthenticated checks if the request has a valid API key
func (m *Middleware) IsAuthenticated(r *http.Request) bool {
	apiKey, found := m.ExtractAPIKey(r)
	if !found {
		return false
	}

	return m.config.IsValidAPIKey(apiKey)
}

// RequireAuth is a middleware that requires authentication
func (m *Middleware) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !m.IsAuthenticated(r) {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("WWW-Authenticate", "Bearer")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error": "unauthorized", "message": "valid API key required in Authorization header"}`))
			return
		}
		next(w, r)
	}
}
