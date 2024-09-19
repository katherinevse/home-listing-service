package handlers

import (
	"app/pkg/auth"
	"encoding/json"
	"net/http"
)

type Handler struct {
	JWTSecretKey string
}

//func NewHandler(jwtSecretKey string) *Handler {
//	return &Handler{JWTSecretKey: cfg
//}

// client, moderator, вернет токен с соответствующим уровнем доступа — обычного пользователя или модератора.
func (h *Handler) DummyLogin(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var requestBody map[string]string
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userType, ok := requestBody["userType"]
	if !ok || (userType != "client" && userType != "moderator") {
		http.Error(w, "Invalid userType", http.StatusBadRequest)
		return
	}

	token, err := auth.GenerateToken(userType, cfg.SecretKey)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	response := map[string]string{"token": token}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
