package auth

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
)

type contextKey string

const UserContextKey contextKey = "user"

func JwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil // TODO: get secret from env
		})
		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Add user's ID to the request context
		ctx := context.WithValue(r.Context(), UserContextKey, claims["user_id"])
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
