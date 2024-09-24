package handler

import (
	"app/internal/dto"
	"app/internal/repository/model"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
	"time"
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

	// Убираем "Bearer " из заголовка
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

	//http.Error(w, "Ok.", http.StatusOK)

	var houseRequest dto.House
	if err := json.NewDecoder(r.Body).Decode(&houseRequest); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	house := model.House{
		City:               houseRequest.City,
		Street:             houseRequest.Street,
		HouseNumber:        houseRequest.HouseNumber,
		YearOfConstruction: houseRequest.YearOfConstruction,
		Developer:          houseRequest.Developer,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	if err := h.houseRepo.CreateHouse(&house); err != nil {
		http.Error(w, "Failed to create house: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := dto.HouseResponse{
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
	json.NewEncoder(w).Encode(response)

}

func (h *Handler) GetFlatsByHouseID(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	//способ получить параметры из URL запроса
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
	u, err := h.tokenManager.ParseJWT(tokenString, h.JWTSecretKey)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		fmt.Printf("Token validation error: %v\n", err)
		return
	}

	flats := make([]model.Flat, 0, 100)

	if u.UserType == "client" {
		flats, err = h.houseRepo.GetApprovedFlatsByHouseID(houseID)
	} else if u.UserType == "moderator" {
		flats, err = h.houseRepo.GetAllFlatsByHouseID(houseID)
	} else {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	if err != nil {
		http.Error(w, "Failed to get flats", http.StatusInternalServerError)
		return
	}

	// Возвращаем результат
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(flats)
}
