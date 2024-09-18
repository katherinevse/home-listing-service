package models

type User struct {
	ID           int    `json:"id" db:"id"`
	Email        string `json:"email" db:"email"` // (уникальный)
	PasswordHash string `json:"-" db:"password_hash"`
	Role         string `json:"role" db:"role"` //  "client" or "moderator"
	//CreatedAt    time.Time `json:"created_at" db:"created_at"` // Дата регистрации
}

//все юзеры отдельная табоица в бд
