package service

import (
	"context"

	"github.com/golovanevvs/gophermart/internal/model"
	"github.com/golovanevvs/gophermart/internal/storage"
	"github.com/golovanevvs/gophermart/internal/transport/http/accrualsystem"
)

type AuthServiceInt interface {
	CreateUser(ctx context.Context, user model.User) (int, error)
	BuildJWTString(ctx context.Context, login, password string) (string, error)
	GetUserIDFromJWT(tokenString string) (int, error)
	GetBalance(ctx context.Context, userID int) (model.Balance, error)
}

type OrderServiceInt interface {
	UploadOrder(ctx context.Context, userID int, orderNumber int) (int, error)
	GetOrders(ctx context.Context, userID int) ([]model.Order, error)
	Withdraw(ctx context.Context, userID int, withdrawOrderNumber string, sum int) error
	Withdrawals(ctx context.Context, userID int) ([]model.Withdrawals, error)
}

type AccrualSystemServiceInt interface {
	GetOrderFromASByOrderNumber(ctx context.Context, orderNumber int) (accrualSystem model.AccrualSystem, err error)
}

type authServiceStr struct {
	st storage.AllStorageInt
}

type orderServiceStr struct {
	st storage.AllStorageInt
}

type accrualSystemServiceStr struct {
	st storage.AllStorageInt
	as accrualsystem.AccrualSystemInt
}

type ServiceStrInt struct {
	AuthServiceInt
	OrderServiceInt
	AccrualSystemServiceInt
}

func NewAuthService(st storage.AllStorageInt) *authServiceStr {
	return &authServiceStr{
		st: st,
	}
}

func NewOrderService(st storage.AllStorageInt) *orderServiceStr {
	return &orderServiceStr{
		st: st,
	}
}

func NewAccrualSystemService(st storage.AllStorageInt, as accrualsystem.AccrualSystemInt) *accrualSystemServiceStr {
	return &accrualSystemServiceStr{
		st: st,
		as: as,
	}
}

func NewService(st *storage.StorageStrInt, as accrualsystem.AccrualSystemInt) *ServiceStrInt {
	return &ServiceStrInt{
		AuthServiceInt:          NewAuthService(st),
		OrderServiceInt:         NewOrderService(st),
		AccrualSystemServiceInt: NewAccrualSystemService(st, as),
	}
}
