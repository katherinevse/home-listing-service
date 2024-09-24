package dto

type Flat struct {
	HouseID    int    `json:"house_id" validate:"required"`
	FlatNumber string `json:"flat_number" validate:"required"`
	Floor      int    `json:"floor" validate:"required"`
	Price      int    `json:"price" validate:"required"`
	RoomsCount int    `json:"rooms_count" validate:"required"`
}

//CREATE TABLE flats (
//id SERIAL PRIMARY KEY,
//house_id INT REFERENCES houses(id) ON DELETE CASCADE,
//flat_number TEXT NOT NULL,
//floor INT NOT NULL,
//price INT NOT NULL,
//rooms_count INT NOT NULL,
//moderation_status TEXT DEFAULT 'created',
//created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
//);
//

type FlatResponse struct {
	HouseID          int    `json:"house_id" validate:"required"`
	FlatNumber       string `json:"flat_number" validate:"required"`
	Floor            int    `json:"floor" validate:"required"`
	Price            int    `json:"price" validate:"required"`
	RoomsCount       int    `json:"rooms_count" validate:"required"`
	ModerationStatus string `json:"moderation_status"` // Статус модерации квартиры (created, approved, declined, on moderation)

}
