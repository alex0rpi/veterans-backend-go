package auth

import (
	"net/http"
	"strings"
)

func (m *JWTManager) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		// Header must be: "Bearer <token>"
		if authHeader == "" {
			http.Error(w, "missing authorization header", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "invalid authorization header", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]
		_, err := m.ValidateToken(tokenString, TokenTypeAccess)
		if err != nil {
			http.Error(w, "invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Pass-on the claims to the next handler
		// OPCIONAL: if we need to access the claimsin within the handler
		// ctx := context.WithValue(r.Context(), "claims", claims)
		// r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
