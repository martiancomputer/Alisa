package api

import (
	"crypto/subtle"
	"net/http"
	"os"
)

// AuthMiddleware inspects the Authorization header against the secure system token.
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Retrieve token from environment variables
		sysToken := os.Getenv("ALISA_DASHBOARD_TOKEN")
		if sysToken == "" {
			http.Error(w, "Internal Server Configuration Defect: Missing Auth Token", http.StatusInternalServerError)
			return
		}

		// Extract token from request header
		authToken := r.Header.Get("Authorization")
		if authToken == "" {
			http.Error(w, "Unauthorized: Authorization Header Missing", http.StatusUnauthorized)
			return
		}

		// Prefix validation tracking (Expected format: "Bearer <token>")
		const prefix = "Bearer "
		if len(authToken) <= len(prefix) || authToken[:len(prefix)] != prefix {
			http.Error(w, "Unauthorized: Invalid Token Formatting", http.StatusUnauthorized)
			return
		}
		tokenValue := authToken[len(prefix):]

		// Execute constant-time comparison to completely neutralize timing attack vectors
		if subtle.ConstantTimeCompare([]byte(tokenValue), []byte(sysToken)) != 1 {
			http.Error(w, "Unauthorized: Token Signatures Do Not Match", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
}