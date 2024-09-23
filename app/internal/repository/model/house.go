package model

import "time"

type House struct {
	ID                 int       `json:"id"`
	City               string    `json:"city" validate:"required"`         // Город
	Street             string    `json:"street" validate:"required"`       // Улица
	HouseNumber        string    `json:"house_number" validate:"required"` // Номер дома
	YearOfConstruction int       `json:"year_built" validate:"required"`
	Developer          *string   `json:"developer,omitempty"` // Необязательное поле
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}
