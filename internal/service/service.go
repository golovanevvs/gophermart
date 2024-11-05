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
	ProcAccrual(userID int, orderNumber int)
	GetAccrual(ctx context.Context, userID int, orderNumber int) (int, error)
	UpdateBalance(ctx context.Context, userID int, accrual int) error
}

type authServiceStr struct {
	st storage.AllStorageInt
}

type orderServiceStr struct {
	st storage.AllStorageInt
	as AccrualSystemServiceInt
}

type accrualSystemServiceStr struct {
	//address string
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

func NewOrderService(st storage.AllStorageInt, accrualServiceAddress string) *orderServiceStr {
	return &orderServiceStr{
		st: st,
		as: NewAccrualSystemService(st, accrualServiceAddress),
	}
}

func NewAccrualSystemService(st storage.AllStorageInt, accrualServiceAddress string) *accrualSystemServiceStr {
	return &accrualSystemServiceStr{
		//address: accrualServiceAddress,
		st: st,
		as: accrualsystem.NewAccrualSystem(accrualServiceAddress),
	}
}

func NewService(st *storage.StorageStrInt, accrualServiceAddress string) *ServiceStrInt {
	return &ServiceStrInt{
		AuthServiceInt:  NewAuthService(st),
		OrderServiceInt: NewOrderService(st, accrualServiceAddress),
	}
}
