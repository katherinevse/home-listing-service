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
