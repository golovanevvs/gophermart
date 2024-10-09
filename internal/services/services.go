package services

import "fmt"

type servicesStr struct {
	db DatabaseInt
}

type DatabaseInt interface {
	SaveToDB() string
}

func New(db DatabaseInt) servicesStr {
	return servicesStr{
		db: db,
	}
}

func (sv servicesStr) RegisterUser() string {
	return fmt.Sprintf("registerUser")
}
