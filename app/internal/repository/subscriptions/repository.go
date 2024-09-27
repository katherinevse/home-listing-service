package subscriptions

import (
	"app/internal/repository/model"
	"context"
	"strconv"
)

type Repo struct {
	db DBPool
}

func NewRepo(db DBPool) *Repo {
	return &Repo{db: db}
}

// TODO context внешне прокинуть
func (r *Repo) CreateSubscription(subscriber *model.Subscription) error {
	query := `INSERT INTO subscriptions (user_id, house_id) VALUES ($1, $2)`
	_, err := r.db.Exec(context.Background(), query, subscriber.UserID, subscriber.HouseID)
	return err
}

func (r *Repo) GetSubscribersByHouseID(houseID int) ([]model.User, error) {
	subscribers := make([]model.User, 0)

	query := `
		SELECT u.id, u.email 
		FROM subscriptions s 
		JOIN users u ON s.user_id = u.id 
		WHERE s.house_id = $1
	`

	rows, err := r.db.Query(context.Background(), query, strconv.Itoa(houseID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user model.User
		if err := rows.Scan(&user.ID, &user.Email); err != nil {
			return nil, err
		}
		subscribers = append(subscribers, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return subscribers, nil
}
