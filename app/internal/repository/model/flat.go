package model

import "time"

type Flat struct {
	ID               int       `json:"id"`                // Уникальный идентификатор квартиры
	HouseID          int       `json:"house_id"`          // Идентификатор дома, к которому относится квартира
	FlatNumber       string    `json:"flat_number"`       // Номер квартиры
	Floor            int       `json:"floor"`             // Этаж, на котором расположена квартира
	Price            int       `json:"price"`             // Цена квартиры
	RoomsCount       int       `json:"rooms_count"`       // Количество комнат
	ModerationStatus string    `json:"moderation_status"` // Статус модерации квартиры (created, approved, declined, on moderation)
	CreatedAt        time.Time `json:"created_at"`        // Дата создания квартиры
}
