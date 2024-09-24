package handler

import (
	"app/internal/repository/model"
)

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
