package model

type User struct {
	UserID   int    `json:"-" db:"user_id"`
	Login    string `json:"login"`
	Password string `json:"password"`
}
