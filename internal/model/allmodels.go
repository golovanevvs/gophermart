package model

import "time"

type User struct {
	UserID       int `json:"-" db:"user_id"`
	Login        string
	Password     string
	PasswordHash string `json:"-" db:"password_hash"`
}

type Order struct {
	OrderID     int       `json:"-" db:"order_id`
	OrderNumber int       `json:"number" db:"order_number"`
	OrderStatus string    `json:"status"`
	UploadedAt  time.Time `json:"uploaded_at" db:"uploaded_at`
}

type Balance struct {
	CurrentPoints int `json:"current" db:"current_points"`
	Withdrawn     int
}

type Accrual struct {
	AccrualPoints int
	AccrualAt     time.Time
}

type WithdrawOrder struct {
	OrderID     int    `json:"-"`
	OrderNumber string `json:"order"`
	Sum         int
	WithdrawAt  time.Time `json:"-"`
}
