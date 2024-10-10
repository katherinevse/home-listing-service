package handler

import (
	"app/internal/dto"
	"app/internal/repository/model"
	"app/pkg/auth"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

//TODO починить айдишники все время 0 в ответе только /house/create, и проверку добавить на одинаковые дома и квартиры в бд

// CreateHouse /house/create
func (h *Handler) CreateHouse(w http.ResponseWriter, r *http.Request) {
	var houseRequest dto.House
	if err := json.NewDecoder(r.Body).Decode(&houseRequest); err != nil {
		h.logger.Error("Invalid request body", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	house := model.House{
		ID:                 houseRequest.ID,
		City:               houseRequest.City,
		Street:             houseRequest.Street,
		HouseNumber:        houseRequest.HouseNumber,
		YearOfConstruction: houseRequest.YearOfConstruction,
		Developer:          houseRequest.Developer,
		// TODO это можно прописать в БД и сделать default now()
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := h.houseRepo.CreateHouse(&house); err != nil {
		h.logger.Error("Failed to create house", "error", err)
		http.Error(w, "Failed to create house: "+err.Error(), http.StatusInternalServerError)
		return

	}

	response := dto.HouseResponse{
		// TODO не по контракту
		ID:                 houseRequest.ID,
		City:               houseRequest.City,
		Street:             houseRequest.Street,
		HouseNumber:        houseRequest.HouseNumber,
		YearOfConstruction: houseRequest.YearOfConstruction,
		Developer:          houseRequest.Developer,
		CreatedAt:          house.CreatedAt.Format(time.RFC3339),
		UpdatedAt:          house.UpdatedAt.Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		h.logger.Error("Failed to encode response", "error", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	h.logger.Info("Successfully created house", "houseID", houseRequest.ID)

}

func (h *Handler) GetFlatsByHouseID(w http.ResponseWriter, r *http.Request) {
	// Получаем параметры из URL запроса
	vars := mux.Vars(r)
	houseID := vars["id"]

	// Извлекаем пользователя из контекста
	user := r.Context().Value("user")
	if user == nil {
		h.logger.Error("User not found in context", "houseID", houseID)
		http.Error(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	u, ok := user.(*auth.Claims)
	if !ok {
		h.logger.Error("Invalid user in context", "houseID", houseID)
		http.Error(w, "Invalid user in context", http.StatusInternalServerError)
		return
	}

	flats := make([]model.Flat, 0, 100)
	var err error

	// Проверяем тип пользователя и выбираем данные в зависимости от его роли
	if u.UserType == "client" {
		h.logger.Info("Fetching approved flats for client", "houseID", houseID)
		flats, err = h.houseRepo.GetApprovedFlatsByHouseID(houseID)
	} else if u.UserType == "moderator" {
		h.logger.Info("Fetching all flats for moderator", "houseID", houseID)
		flats, err = h.houseRepo.GetAllFlatsByHouseID(houseID)
	} else {
		h.logger.Error("Access denied", "userType", u.UserType, "houseID", houseID)
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	if err != nil {
		h.logger.Error("Failed to get flats", "houseID", houseID, "error", err)
		http.Error(w, "Failed to get flats", http.StatusInternalServerError)
		return
	}

	// Возвращаем результат
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(flats)
	if err != nil {
		h.logger.Error("Failed to encode response", "houseID", houseID, "error", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
	h.logger.Info("Successfully fetched flats", "houseID", houseID, "count", len(flats))

}
