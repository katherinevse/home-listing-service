package model

type Subscription struct {
	ID      int    `json:"id"`
	UserID  int    `json:"user_id"`
	HouseID string `json:"house_id"`
}
