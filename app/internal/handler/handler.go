package handler

import (
	"app/pkg/auth"
	"os"
)

type Handler struct {
	JWTSecretKey string
	tokenManager auth.TokenManager
	userRepo     UserRepository
	houseRepo    HouseRepository
	flatRepo     FlatRepository
}

func New(tokenManager auth.TokenManager, userRepo UserRepository, houseRepo HouseRepository, flatRepo FlatRepository) *Handler {
	return &Handler{JWTSecretKey: os.Getenv("JWT_SECRET_KEY"), tokenManager: tokenManager, userRepo: userRepo, houseRepo: houseRepo, flatRepo: flatRepo}
}
