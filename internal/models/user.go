package models

// User данные пользователя
type User struct {
	ID           int    `db:"id"`
	Name         string `db:"name"`
	PasswordHash []byte `db:"password"`
}

// Credentials учетные данные
type Credentials struct {
	Username string `json:"username" example:"user1"`
	Password string `json:"password" example:"password123"`
}
