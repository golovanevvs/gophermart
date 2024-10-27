package model

import "time"

type Accrual struct {
	AccrualPoints int
	AccrualAt     time.Time
}
