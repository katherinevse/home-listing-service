package handler

import (
	"app/internal/repository/model"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

// CreateSubscription  /house/{id}/subscribe
func (h *Handler) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	houseID := vars["id"]
	u, ok := r.Context().Value("user").(*model.User)
	if !ok {
		http.Error(w, "Failed to retrieve user from context", http.StatusInternalServerError)
		return
	}

	subscription := model.Subscription{
		UserID:  u.ID,
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
