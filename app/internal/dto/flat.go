package dto

type Flat struct {
	HouseID    int `json:"house_id" validate:"required"`
	FlatNumber int `json:"flat_number" validate:"required"`
	Floor      int `json:"floor" validate:"required"`
	Price      int `json:"price" validate:"required"`
	RoomsCount int `json:"rooms_count" validate:"required"`
}

type FlatResponse struct {
	HouseID          int    `json:"house_id" validate:"required"`
	FlatNumber       int    `json:"flat_number" validate:"required"`
	Floor            int    `json:"floor" validate:"required"`
	Price            int    `json:"price" validate:"required"`
	RoomsCount       int    `json:"rooms_count" validate:"required"`
	ModerationStatus string `json:"moderation_status"` // (created, approved, declined, on moderation)

}
