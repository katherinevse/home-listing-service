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

// TODO ctx передать снаружи
func (r *Repo) CreateFlat(flat *model.Flat) error {
	query := `INSERT INTO flats (house_id, flat_number, floor, price, rooms_count, moderation_status, created_at) 
              VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	return r.db.QueryRow(context.Background(),
		query,
		flat.HouseID,
		flat.FlatNumber,
		flat.Floor,
		flat.Price,
		flat.RoomsCount,
		flat.ModerationStatus,
		flat.CreatedAt).Scan(&flat.ID)
}

func (r *Repo) GetFlatsOnModeration() ([]model.Flat, error) {
	query := `SELECT id, house_id, flat_number, floor, price, rooms_count, moderation_status, created_at 
			  FROM flats 
			  WHERE moderation_status = 'on moderation'`

	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	flats := make([]model.Flat, 0)

	for rows.Next() {
		var flat model.Flat
		err := rows.Scan(&flat.ID,
			&flat.HouseID,
			&flat.FlatNumber,
			&flat.Floor,
			&flat.Price,
			&flat.RoomsCount,
			&flat.ModerationStatus,
			&flat.CreatedAt)

		if err != nil {
			return nil, err
		}

		flats = append(flats, flat)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return flats, nil
}

//контексты
//логгер
