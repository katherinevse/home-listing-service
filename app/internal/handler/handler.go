package handler

import (
	"app/pkg/auth"
	"os"
)

type Handler struct {
	JWTSecretKey string
	userRepo     UserRepository
	tokenManager auth.TokenManager
	houseRepo    HouseRepository
}

func New(tokenManager auth.TokenManager, userRepo UserRepository, houseRepo HouseRepository) *Handler {
	return &Handler{JWTSecretKey: os.Getenv("JWT_SECRET_KEY"), tokenManager: tokenManager, userRepo: userRepo, houseRepo: houseRepo}
}
