package handler

import (
	"app/internal/repository/model"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strings"
)

// CreateSubscription  /house/{id}/subscribe
func (h *Handler) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	houseID := vars["id"]

	authHeader := r.Header.Get("Authorization")
	bearerPrefix := "Bearer "
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		http.Error(w, "Missing or invalid Authorization header", http.StatusBadRequest)
		return
	}

	// Парсинг токена
	tokenString := strings.TrimPrefix(authHeader, bearerPrefix)
	u, err := h.tokenManager.ParseJWT(tokenString)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		fmt.Printf("Token validation error: %v\n", err)
		return
	}

	subscription := model.Subscription{
		UserID:  u.UserID,
		HouseID: houseID,
	}

	if err := h.subscriptionRepo.CreateSubscription(&subscription); err != nil {
		http.Error(w, "Failed to create subscription", http.StatusInternalServerError)
		log.Printf("Failed to create subscription: %v", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	response, err := w.Write([]byte("Subscription created successfully"))

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Printf("Failed to encode response: %v", err)
		return
	}
}
