package middleware

import (
	"net/http"
	"strings"
)

//type APPHandler func(w http.ResponseWriter, r *http.Request) error

func AuthMiddleware(next http.HandlerFunc, tokenManager TokenManager) http.HandlerFunc {
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
			return
		}

		if u.UserType != "moderator" {
			http.Error(w, "Access denied. Only moderators can perform this action.", http.StatusForbidden)
			return
		}

		// Если все в порядке, продолжаем вызов следующего обработчика
		next(w, r)
	}
}
