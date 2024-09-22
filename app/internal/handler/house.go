package handler

import (
	"fmt"
	"net/http"
	"strings"
)

// CreateHouse /house/create
func (h *Handler) CreateHouse(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	authHeader := r.Header.Get("Authorization")
	bearerPrefix := "Bearer "

	// Проверяем наличие заголовка и его корректность
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		http.Error(w, "Missing or invalid Authorization header", http.StatusBadRequest)
		return
	}

	// Убираем префикс "Bearer " из заголовка
	tokenString := strings.TrimPrefix(authHeader, bearerPrefix)

	u, err := h.tokenManager.ParseJWT(tokenString, h.JWTSecretKey)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		fmt.Printf("Token validation error: %v\n", err)

	}
	if u.UserType == "client" {
		http.Error(w, "Access denied. Only moderators can perform this action.", http.StatusForbidden)
		fmt.Println("Attempt to access moderator-only endpoint by non-moderator user -->", u.Email, u.UserID, u.UserType)
		return
	}

	http.Error(w, "Ok.", http.StatusOK)

}
