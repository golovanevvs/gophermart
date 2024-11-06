package model

import "time"

type User struct {
	ID           int    `json:"-"`
	Login        string `json:"login"`
	Password     string `json:"password"`
	PasswordHash string `json:"-"`
}

type Order struct {
	ID         int       `json:"-"`
	Number     int       `json:"number"`
	Status     string    `json:"status"`
	Accrual    float64   `json:"accrual,omitempty"`
	UploadedAt time.Time `json:"uploaded_at"`
}

type Balance struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

type Withdrawals struct {
	ID             int       `json:"-"`
	NewOrderNumber string    `json:"order"`
	Sum            float64   `json:"sum"`
	ProcessedAt    time.Time `json:"processed_at"`
}

type AccrualSystem struct {
	OrderNumber string  `json:"order"`
	Status      string  `json:"status"`
	Accrual     float64 `json:"accrual"`
	RetryAfter  int     `json:"-"`
	Message     string  `json:"-"`
}
