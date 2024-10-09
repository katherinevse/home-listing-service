package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

// Auth - middleware для аутентификации и проверки роли
func Auth(next http.HandlerFunc, tokenManager TokenManager, moderatorCheck bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		bearerPrefix := "Bearer "

		if !strings.HasPrefix(authHeader, bearerPrefix) {
			http.Error(w, "Missing or invalid Authorization header", http.StatusBadRequest)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, bearerPrefix)

		u, err := tokenManager.ParseJWT(tokenString)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			fmt.Printf("Token validation error: %v\n", err)
			return
		}

		if moderatorCheck && u.UserType != "moderator" {
			http.Error(w, "Access denied. Only moderators can perform this action.", http.StatusForbidden)
			fmt.Println("Attempt to access moderator-only endpoint by non-moderator user -->", u.Email, u.UserID, u.UserType)
			return
		}

		// пользователя в контекст запроса для дома
		ctx := context.WithValue(r.Context(), "user", u)
		next(w, r.WithContext(ctx))
	}
}
