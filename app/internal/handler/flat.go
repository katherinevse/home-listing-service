package handler

import (
	"app/internal/dto"
	"app/internal/repository/model"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	maxPrice = 1000000
	maxRooms = 5
)

func (h *Handler) CreateFlat(w http.ResponseWriter, r *http.Request) {
	var flatRequest dto.Flat
	if err := json.NewDecoder(r.Body).Decode(&flatRequest); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Если значения в пределах, ставим статус "одобрено" иначе на модерацию
	var moderationStatus string
	if flatRequest.Price > maxPrice || flatRequest.RoomsCount > maxRooms {
		moderationStatus = "on moderation"
	} else {
		moderationStatus = "approved"
	}

	flat := model.Flat{
		HouseID:          flatRequest.HouseID,
		FlatNumber:       flatRequest.FlatNumber,
		Price:            flatRequest.Price,
		RoomsCount:       flatRequest.RoomsCount,
		ModerationStatus: moderationStatus,
		CreatedAt:        time.Now(),
	}

	if err := h.flatRepo.CreateFlat(&flat); err != nil {
		http.Error(w, "Failed to create flat", http.StatusInternalServerError)
		return
	}

	//TODO горутинку добавить стоит
	message := "New flat created!"
	if err := h.producer.PublishNotification(flat.HouseID, flat.FlatNumber, message); err != nil {
		log.Printf("Failed to send Kafka notification for houseID %s, flatNumber %s: %v", flat.HouseID, flat.FlatNumber, err)
		// TODO продолжать ли , даже если отправка уведомления не удалась
	}

	response := dto.FlatResponse{
		HouseID:          flat.HouseID,
		FlatNumber:       flat.FlatNumber,
		Price:            flat.Price,
		RoomsCount:       flat.RoomsCount,
		ModerationStatus: flat.ModerationStatus,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) GetModerationFlats(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	bearerPrefix := "Bearer "

	if !strings.HasPrefix(authHeader, bearerPrefix) {
		http.Error(w, "Missing or invalid Authorization header", http.StatusBadRequest)
		return
	}

	tokenString := strings.TrimPrefix(authHeader, bearerPrefix)

	u, err := h.tokenManager.ParseJWT(tokenString, h.JWTSecretKey)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		// TODO logger
		fmt.Printf("Token validation error: %v\n", err)
		return
	}

	if u.UserType == "client" {
		http.Error(w, "Access denied. Only moderators can perform this action.", http.StatusForbidden)
		fmt.Println("Attempt to access moderator-only endpoint by non-moderator user -->", u.Email, u.UserID, u.UserType)
		return
	}

	flats, err := h.flatRepo.GetFlatsOnModeration()
	if err != nil {
		http.Error(w, "Failed to get flats", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(flats)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Printf("Failed to encode response: %v", err)
		return
	}

}
