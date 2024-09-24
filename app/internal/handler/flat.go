package handler

import (
	"app/internal/dto"
	"app/internal/repository/model"
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

func (h *Handler) CreateFlat(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	authHeader := r.Header.Get("Authorization")
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	_, err := h.tokenManager.ParseJWT(tokenString, h.JWTSecretKey)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	var flatRequest dto.Flat
	if err := json.NewDecoder(r.Body).Decode(&flatRequest); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	flat := model.Flat{
		HouseID:          flatRequest.HouseID,
		FlatNumber:       flatRequest.FlatNumber,
		Floor:            flatRequest.Floor,
		Price:            flatRequest.Price,
		RoomsCount:       flatRequest.RoomsCount,
		ModerationStatus: "created, waiting for moderation",
		CreatedAt:        time.Now(),
	}

	if err := h.flatRepo.CreateFlat(&flat); err != nil {
		http.Error(w, "Failed to create flat", http.StatusInternalServerError)
		return
	}

	flatResponse := dto.FlatResponse{
		HouseID:          flat.HouseID,
		FlatNumber:       flat.FlatNumber,
		Floor:            flat.Floor,
		Price:            flat.Price,
		RoomsCount:       flat.RoomsCount,
		ModerationStatus: flat.ModerationStatus,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(flatResponse)
}
