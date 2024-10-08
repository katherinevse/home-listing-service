package handler

import (
	"app/internal/dto"
	"app/internal/repository/model"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"io"
	"net/http"
	"time"
)

// Register /register
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var user dto.Register
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		h.logger.Error("Invalid request payload", "error", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			h.logger.Error("Error closing request body", "error", err)
		}
	}(r.Body)

	if user.Email == "" || user.Password == "" || user.UserType == "" {
		h.logger.Error("Invalid user data", "user", user)
		http.Error(w, "Invalid user data", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		h.logger.Error("Error hashing password", "error", err)
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	u := model.User{
		Email:    user.Email,
		Password: string(hashedPassword),
		UserType: user.UserType,
	}

	err = h.userRepo.CreateUser(&u, hashedPassword)
	if err != nil {
		h.logger.Error("Error saving user to database", "error", err)
		http.Error(w, fmt.Sprintf("Error saving user to database: %s", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte("User registered successfully"))
	if err != nil {
		h.logger.Error("Error writing response", "error", err)
		http.Error(w, "Error writing response", http.StatusInternalServerError)
		return
	}

	h.logger.Info("User registered successfully", "email", user.Email)
}

// Login /login
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var user dto.Login
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		h.logger.Error("Invalid request payload", "error", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			h.logger.Error("Error closing request body", "error", err)
		}
	}(r.Body)

	if user.Email == "" || user.Password == "" {
		h.logger.Error("Invalid user data", "user", user)
		http.Error(w, "Invalid user data", http.StatusBadRequest)
		return
	}

	dbUser, err := h.userRepo.GetUserByEmail(user.Email)
	if err != nil {
		h.logger.Error("Invalid email or password", "email", user.Email)
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password)); err != nil {
		h.logger.Error("Invalid email or password", "email", user.Email)
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	tokenString, err := h.tokenManager.GenerateJWT(dbUser.ID, dbUser.UserType)
	if err != nil {
		h.logger.Error("Failed to generate JWT token", "error", err)
		http.Error(w, "Failed to generate JWT token", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "jwt_token",
		Value:    tokenString,
		Expires:  time.Now().Add(24 * time.Hour),
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
		HttpOnly: true,
		Secure:   true,
	})

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("Login successful"))
	if err != nil {
		h.logger.Error("Error writing response", "error", err)
		http.Error(w, "Error writing response", http.StatusInternalServerError)
		return
	}

	h.logger.Info("User logged in successfully", "email", user.Email)

}
