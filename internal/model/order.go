package model

import "time"

type Order struct {
	OrderID     int
	OrderNumber int
	OrderStatus string
	UploadedAt  time.Time
}
