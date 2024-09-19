package user

import (
	"app/internal/model"
	"context"
)

type Repo struct {
	db DBPool
}

func NewRepo(db DBPool) *Repo {
	return &Repo{db: db}
}

func (r *Repo) CreateUser(user *model.User, hashedPassword []byte) error {
	query := `INSERT INTO users (email, password, usertype) VALUES ($1, $2, $3)`
	_, err := r.db.Exec(context.Background(), query, user.Email, hashedPassword, user.UserType)
	if err != nil {
		return err
	}
	return nil
}
