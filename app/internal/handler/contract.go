package handler

import (
	"app/internal/model"
)

type UserRepository interface {
	CreateUser(user *model.User, hashedPassword []byte) error
}
