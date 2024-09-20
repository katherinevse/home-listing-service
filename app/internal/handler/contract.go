package handler

import (
	"app/internal/repository/model"
)

type UserRepository interface {
	CreateUser(user *model.User, hashedPassword []byte) error
	GetUserByEmail(email string) (*model.User, error)
}
