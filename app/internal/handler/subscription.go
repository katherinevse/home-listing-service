package handler

import (
	"app/internal/repository/model"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

// CreateSubscription  /house/{id}/subscribe
func (h *Handler) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	houseID := vars["id"]
	u, ok := r.Context().Value("user").(*model.User)
	if !ok {
		h.logger.Error("Failed to retrieve user from context", "houseID", houseID)
		http.Error(w, "Failed to retrieve user from context", http.StatusInternalServerError)
		return
	}

	subscription := model.Subscription{
		UserID:  u.ID,
		HouseID: houseID,
	}

	if err := h.subscriptionRepo.CreateSubscription(&subscription); err != nil {
		h.logger.Error("Failed to create subscription", "houseID", houseID, "userID", u.ID, "error", err)
		http.Error(w, "Failed to create subscription", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	response, err := w.Write([]byte("Subscription created successfully"))

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		h.logger.Error("Failed to encode response", "houseID", houseID, "userID", u.ID, "error", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
	h.logger.Info("Successfully created subscription", "houseID", houseID, "userID", u.ID)

}
