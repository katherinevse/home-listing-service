package dto

type Register struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	UserType string `json:"userType"`
}
