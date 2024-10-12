package service

import (
	"context"

	"github.com/golovanevvs/gophermart/internal/model"
	"github.com/golovanevvs/gophermart/internal/storage"
)

type AuthInt interface {
	CreateUser(ctx context.Context, user model.User) error
}

type ServiceStrInt struct {
	AuthInt
}

func NewService(st *storage.StorageStrInt) *ServiceStrInt {
	return &ServiceStrInt{
		AuthInt: NewAuthService(st),
	}
}
