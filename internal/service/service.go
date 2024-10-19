package service

import (
	"context"

	"github.com/golovanevvs/gophermart/internal/customerrors"
	"github.com/golovanevvs/gophermart/internal/model"
	"github.com/golovanevvs/gophermart/internal/storage"
)

type AuthInt interface {
	CreateUser(ctx context.Context, user model.User) (int, error)
	BuildJWTString(ctx context.Context, login, password string) (string, customerrors.CustomError)
	GetUserIDFromJWT(tokenString string) (int, error)
}

type authServiceStr struct {
	st storage.AllStorageInt
}

type ServiceStrInt struct {
	AuthInt
}

func NewAuthService(st storage.AllStorageInt) *authServiceStr {
	return &authServiceStr{
		st: st,
	}
}

func NewService(st *storage.StorageStrInt) *ServiceStrInt {
	return &ServiceStrInt{
		AuthInt: NewAuthService(st),
	}
}
