package handler

import (
	"app/internal/dto"
	"app/internal/repository/model"
	"encoding/json"
	"net/http"
	"time"
)

const (
	maxPrice = 1000000
	maxRooms = 5
)

func (h *Handler) CreateFlat(w http.ResponseWriter, r *http.Request) {

	var flatRequest dto.Flat
	if err := json.NewDecoder(r.Body).Decode(&flatRequest); err != nil {
		h.logger.Error("Invalid request body", "error", err)
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
		h.logger.Error("Failed to create flat", "error", err)
		http.Error(w, "Failed to create flat", http.StatusInternalServerError)
		return
	}

	go func() {
		message := "New flat created!"
		if err := h.producer.PublishNotification(flat.HouseID, flat.FlatNumber, message); err != nil {
			h.logger.Error("Failed to send Kafka notification", "houseID", flat.HouseID, "flatNumber", flat.FlatNumber, "error", err)
		}
	}()

	response := dto.FlatResponse{
		HouseID:          flat.HouseID,
		FlatNumber:       flat.FlatNumber,
		Price:            flat.Price,
		RoomsCount:       flat.RoomsCount,
		ModerationStatus: flat.ModerationStatus,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		h.logger.Error("Failed to encode response", "error", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
	return
}

func (h *Handler) GetModerationFlats(w http.ResponseWriter, r *http.Request) {
	flats, err := h.flatRepo.GetFlatsOnModeration()
	if err != nil {
		h.logger.Error("Failed to get flats on moderation", "error", err)
		http.Error(w, "Failed to get flats", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(flats)
	if err != nil {
		h.logger.Error("Failed to encode response", "error", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	h.logger.Info("Successfully returned flats on moderation")

}
