package dto

import "time"

// CreateHouse - структура для создания дома
type House struct {
	ID                 int       `json:"id"`
	City               string    `json:"city" validate:"required"`         // Город
	Street             string    `json:"street" validate:"required"`       // Улица
	HouseNumber        string    `json:"house_number" validate:"required"` // Номер дома
	YearOfConstruction int       `json:"year_of_construction" validate:"required"`
	Developer          *string   `json:"developer,omitempty"` // Необязательное поле
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// HouseResponse - структура для ответа с данными дома
type HouseResponse struct {
	ID                 int     `json:"id"`                   // Уникальный идентификатор дома
	City               string  `json:"city"`                 // Город
	Street             string  `json:"street"`               // Улица
	HouseNumber        string  `json:"house_number"`         // Номер дома
	YearOfConstruction int     `json:"year_of_construction"` // Год постройки
	Developer          *string `json:"developer,omitempty"`  // Необязательное поле
	CreatedAt          string  `json:"created_at"`           // Дата создания дома
	UpdatedAt          string  `json:"updated_at"`           // Дата последнего обновления
}
