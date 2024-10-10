package service

import "github.com/golovanevvs/gophermart/internal/storage"

type AuthInt interface {
	CreateUser() string
}

type ServiceStr struct {
	AuthInt
}

func NewService(st *storage.StorageStr) *ServiceStr {
	return &ServiceStr{
		AuthInt: NewAuthService(st),
	}
}
