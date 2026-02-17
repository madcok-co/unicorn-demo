package middleware

import (
	"net/http"
	"strings"
)

// AuthMiddleware checks for authentication token
// In a real app, this would validate JWT token
func AuthMiddleware() func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// Skip auth for public endpoints
			if isPublicEndpoint(r.URL.Path) {
				next(w, r)
				return
			}

			// Get token from Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Unauthorized: missing token", http.StatusUnauthorized)
				return
			}

			// Extract bearer token
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "Unauthorized: invalid token format", http.StatusUnauthorized)
				return
			}

			token := parts[1]

			// In a real app, validate JWT and extract user ID
			// For demo, we'll extract user ID from token
			userID := extractUserIDFromToken(token)
			if userID == "" {
				http.Error(w, "Unauthorized: invalid token", http.StatusUnauthorized)
				return
			}

			// Store user ID in header (will be accessible in handler via ctx.Get("user_id"))
			r.Header.Set("X-User-ID", userID)

			next(w, r)
		}
	}
}

// isPublicEndpoint checks if endpoint doesn't require authentication
func isPublicEndpoint(path string) bool {
	publicPaths := []string{
		"/health",
		"/api/v1/auth/login",
		"/api/v1/auth/callback",
		"/api/v2/auth/login",
		"/api/v2/auth/callback",
	}

	for _, p := range publicPaths {
		if path == p || strings.HasPrefix(path, p) {
			return true
		}
	}

	return false
}

// extractUserIDFromToken extracts user ID from token
// In a real app, this would decode and validate JWT
func extractUserIDFromToken(token string) string {
	// Mock implementation for demo
	if strings.HasPrefix(token, "jwt-token-") {
		return strings.TrimPrefix(token, "jwt-token-")
	}

	// For OAuth2 tokens, we'd validate with the provider
	// For now, return a default user
	return "user-dev-acme"
}
