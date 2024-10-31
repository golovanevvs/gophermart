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
	UploadedAt time.Time `json:"uploaded_at"`
}

type Balance struct {
	Current   int
	Withdrawn int
}

type Accrual struct {
	ID            int
	AccrualPoints int
	AccrualAt     time.Time
}

type Withdrawals struct {
	ID          int    `json:"-"`
	OrderNumber string `json:"order"`
	Sum         int
	WithdrawAt  time.Time `json:"-"`
}
