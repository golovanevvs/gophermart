package model

import "time"

type User struct {
	ID           int `json:"-"`
	Login        string
	Password     string
	PasswordHash string `json:"-"`
}

type Order struct {
	ID         int `json:"-"`
	Number     int
	Status     string
	Accrual    int
	UploadedAt time.Time `json:"uploaded_at"`
}

type Balance struct {
	Current   int
	Withdrawn int
}

type Withdrawals struct {
	ID             int    `json:"-"`
	NewOrderNumber string `json:"order"`
	Sum            int
	ProcessedAt    time.Time `json:"processed_at"`
}

type AccrualSystem struct {
	OrderNumber string `json:"order"`
	Status      string
	Accrual     int
	RetryAfter  int    `json:"-"`
	Message     string `json:"-"`
}
