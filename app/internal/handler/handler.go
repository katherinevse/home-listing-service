package handler

import (
	"os"
)

type Handler struct {
	JWTSecretKey     string
	tokenManager     TokenManager
	userRepo         UserRepository
	houseRepo        HouseRepository
	flatRepo         FlatRepository
	subscriptionRepo SubscriptionRepository
	producer         ProducerManager
	consumer         ConsumerManager
}

func New(
	tokenManager TokenManager,
	userRepo UserRepository,
	houseRepo HouseRepository,
	flatRepo FlatRepository,
	subscriptionRepo SubscriptionRepository,
	producer ProducerManager,
	consumer ConsumerManager,
) *Handler {
	return &Handler{
		JWTSecretKey:     os.Getenv("JWT_SECRET_KEY"),
		tokenManager:     tokenManager,
		userRepo:         userRepo,
		houseRepo:        houseRepo,
		flatRepo:         flatRepo,
		subscriptionRepo: subscriptionRepo,
		producer:         producer,
		consumer:         consumer,
	}
}
