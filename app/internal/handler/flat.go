package handler

import (
	"app/internal/dto"
	"app/internal/repository/model"
	"encoding/json"
	"net/http"
	"time"
)

func (h *Handler) CreateFlat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var flatRequest dto.Flat
	if err := json.NewDecoder(r.Body).Decode(&flatRequest); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	const maxPrice = 100000000
	const maxRooms = 10

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
