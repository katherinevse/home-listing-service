package house

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

func (r *Repo) CreateHouse(house *model.House) error {
	query := `INSERT INTO houses (city, street, house_number, year_built, developer, created_at, updated_at) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.Exec(context.Background(), query,
		house.City,
		house.Street,
		house.HouseNumber,
		house.YearOfConstruction,
		house.Developer,
		house.CreatedAt,
		house.UpdatedAt)

	if err != nil {
		return err
	}
	return nil
}

// TODO лежит в house
func (r *Repo) GetAllFlatsByHouseID(houseID string) ([]model.Flat, error) {
	flats := make([]model.Flat, 0)
	query := `SELECT id, house_id, flat_number, floor, price, rooms_count, moderation_status, created_at
			  FROM flats
			  WHERE house_id = $1`

	rows, err := r.db.Query(context.Background(), query, houseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var flat model.Flat
		err = rows.Scan(
			&flat.ID,
			&flat.HouseID,
			&flat.FlatNumber,
			&flat.Floor,
			&flat.Price,
			&flat.RoomsCount,
			&flat.ModerationStatus,
			&flat.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		flats = append(flats, flat)
	}

	return flats, nil
}

func (r *Repo) GetApprovedFlatsByHouseID(houseID string) ([]model.Flat, error) {
	flats := make([]model.Flat, 0, 100)
	query := `SELECT id, house_id, flat_number, floor, price, rooms_count, moderation_status, created_at
			  FROM flats
			  WHERE house_id = $1 AND moderation_status = 'approved'`

	rows, err := r.db.Query(context.Background(), query, houseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var flat model.Flat
		err := rows.Scan(
			&flat.ID,
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

	return flats, nil
}
