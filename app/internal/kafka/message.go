package kafka

// TODO модель сообщений или сообщение все таки?
type NotificationMessage struct {
	HouseID    int    `json:"house_id"`
	FlatNumber int    `json:"flat_number"`
	Message    string `json:"message"`
}
