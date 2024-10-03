package handler

import (
	"app/internal/repository/model"
	"app/pkg/auth"
)

//go:generate mockgen -source=contract.go -destination=mocks/mockTokenManager.go
type TokenManager interface {
	GenerateJWT(userID int, userType string) (string, error)
	ParseJWT(tokenStr string) (*auth.Claims, error)
}

type UserRepository interface {
	CreateUser(user *model.User, hashedPassword []byte) error
	GetUserByEmail(email string) (*model.User, error)
}

type HouseRepository interface {
	CreateHouse(house *model.House) error
	GetAllFlatsByHouseID(houseID string) ([]model.Flat, error)
	GetApprovedFlatsByHouseID(houseID string) ([]model.Flat, error)
}

type FlatRepository interface {
	CreateFlat(flat *model.Flat) error
	GetFlatsOnModeration() ([]model.Flat, error)
}

type SubscriptionRepository interface {
	CreateSubscription(subscriber *model.Subscription) error
	GetSubscribersByHouseID(houseID int) ([]model.User, error)
}

type ProducerManager interface {
	PublishNotification(houseID int, flatNumber int, message string) error
}

type ConsumerManager interface {
	Listen(topic string)
}
