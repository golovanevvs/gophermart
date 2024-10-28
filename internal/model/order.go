package model

import "time"

type Order struct {
	OrderID     int       `json:"-"`
	OrderNumber int       `json:"number"`
	OrderStatus string    `json:"status"`
	UploadedAt  time.Time `json:"uploaded_at"`
}
