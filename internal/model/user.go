package model

type User struct {
	UserID       int    `json:"-" db:"user_id"`
	Login        string `json:"login"`
	Password     string `json:"password"`
	PasswordHash string `db:"password_hash"`
	Points       int    `db:"points"`
}
