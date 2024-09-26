package handler

import (
	"os"
)

type Handler struct {
	JWTSecretKey string
	tokenManager TokenManager
	userRepo     UserRepository
	houseRepo    HouseRepository
	flatRepo     FlatRepository
	producer     ProducerManager
}

func New(
	tokenManager TokenManager,
	userRepo UserRepository,
	houseRepo HouseRepository,
	flatRepo FlatRepository,
	producer ProducerManager,
) *Handler {
	return &Handler{
		JWTSecretKey: os.Getenv("JWT_SECRET_KEY"),
		tokenManager: tokenManager,
		userRepo:     userRepo,
		houseRepo:    houseRepo,
		flatRepo:     flatRepo,
		producer:     producer,
	}
}
