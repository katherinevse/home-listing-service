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

func (r *Repo) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	query := `SELECT id, email, password, usertype FROM users WHERE email=$1`
	err := r.db.QueryRow(context.Background(), query, email).Scan(&user.ID, &user.Email, &user.Password, &user.UserType)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
