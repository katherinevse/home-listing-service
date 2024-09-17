package models

import (
	"time"
)

type House struct {
	ID               int       `json:"id" db:"id"` // Уникальный
	Address          string    `json:"address" db:"address"`
	ConstructionYear int       `json:"construction_year" db:"construction_year"` // Год постройки дома
	Developer        *string   `json:"developer,omitempty" db:"developer"`       // Застройщик (может быть NULL)
	CreatedAt        time.Time `json:"created_at" db:"created_at"`               // Дата добавления дома в базу данных
	LastFlatAdded    time.Time `json:"last_flat_added" db:"last_flat_added"`     // Дата последнего добавления квартиры
}

// Таблица для домов — "houses".
