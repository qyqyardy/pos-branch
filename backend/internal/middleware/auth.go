package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const (
	UserIDKey contextKey = "user_id"
	RoleKey   contextKey = "role"
)

func JWT(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		auth := r.Header.Get("Authorization")
		if auth == "" {
			http.Error(w, "Missing token", 401)
			return
		}

		parts := strings.Fields(auth)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			http.Error(w, "Invalid token", 401)
			return
		}

		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			http.Error(w, "Server misconfigured", 500)
			return
		}

		token, err := jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
			if token.Method == nil || token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", 401)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Invalid token", 401)
			return
		}

		userID, ok := claims["user_id"].(string)
		if !ok || userID == "" {
			http.Error(w, "Invalid token", 401)
			return
		}

		role, ok := claims["role"].(string)
		if !ok || role == "" {
			http.Error(w, "Invalid token", 401)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		ctx = context.WithValue(ctx, RoleKey, role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
