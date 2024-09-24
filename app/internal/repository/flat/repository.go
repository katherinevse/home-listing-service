package flat

import (
	"app/internal/repository/model"
	"context"
)

type Repo struct {
	db DBPool
}

func NewRepo(db DBPool) *Repo {
	return &Repo{db: db}
}

func (r *Repo) CreateFlat(flat *model.Flat) error {
	query := `INSERT INTO flats (house_id, flat_number, floor, price, rooms_count, moderation_status, created_at) 
              VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	return r.db.QueryRow(context.Background(), query, flat.HouseID, flat.FlatNumber, flat.Floor, flat.Price, flat.RoomsCount, flat.ModerationStatus, flat.CreatedAt).Scan(&flat.ID)
}
