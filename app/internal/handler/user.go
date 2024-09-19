package handler

import (
	"app/internal/model"
	"app/pkg/auth"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"io"
	"log"
	"net/http"
	"os"
)

type Handler struct {
	JWTSecretKey string
	userRepo     UserRepository
	tokenManager auth.TokenManager
}

func New(tokenManager auth.TokenManager, userRepo UserRepository) *Handler {
	return &Handler{JWTSecretKey: os.Getenv("JWT_SECRET_KEY"), tokenManager: tokenManager, userRepo: userRepo}
}

// Register /register
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Error closing request body: %v", err)
		}
	}(r.Body)

	if user.Email == "" || user.Password == "" || user.UserType == "" {
		http.Error(w, "Invalid user data", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	err = h.userRepo.CreateUser(&user, hashedPassword)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error saving user to database: %s", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte("User registered successfully"))
	if err != nil {
		log.Printf("Error writing response: %v", err)
		return
	}
}
