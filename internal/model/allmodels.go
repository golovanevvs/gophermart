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
	Accrual    int       `json:"accrual,omitempty"`
	UploadedAt time.Time `json:"uploaded_at"`
}

type Balance struct {
	Current   int `json:"current"`
	Withdrawn int `json:"withdrawn"`
}

type Withdrawals struct {
	ID             int       `json:"-"`
	NewOrderNumber string    `json:"order"`
	Sum            int       `json:"sum"`
	ProcessedAt    time.Time `json:"processed_at"`
}

type AccrualSystem struct {
	OrderNumber string `json:"order"`
	Status      string `json:"status"`
	Accrual     int    `json:"accrual"`
	RetryAfter  int    `json:"-"`
	Message     string `json:"-"`
}
