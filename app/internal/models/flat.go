package models

import (
	"time"
)

type Flat struct {
	ID         int       `json:"id" db:"id"`
	HouseID    int       `json:"house_id" db:"house_id"` // ID дома, к которому привязана квартира
	FlatNumber int       `json:"flat_number" db:"flat_number"`
	Price      int       `json:"price" db:"price"`
	Rooms      int       `json:"rooms" db:"rooms"`           // Количество комнат
	Status     string    `json:"status" db:"status"`         // Статус модерации: "created", "on moderation", "approved", "declined"
	CreatedAt  time.Time `json:"created_at" db:"created_at"` // Дата создания квартиры
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"` // Дата последнего изменения статуса
}

// Таблица для квартир — "flats".
