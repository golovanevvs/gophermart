package storage

import (
	"context"

	"github.com/golovanevvs/gophermart/internal/model"
	"github.com/golovanevvs/gophermart/internal/storage/postgres"
	"github.com/jmoiron/sqlx"
)

type AllStorageInt interface {
	SaveUser(ctx context.Context, user model.User) (int, error)
	LoadUserByLoginPasswordHash(ctx context.Context, login, password string) (model.User, error)
	SaveOrderNumberByUserID(ctx context.Context, userID int, orderNumber int) (int, error)
	LoadUserIDByOrderNumber(ctx context.Context, orderNumber int) (int, error)
	LoadOrderByUserID(ctx context.Context, userID int) ([]model.Order, error)
	LoadBalanceByUserID(ctx context.Context, userID int) (model.Balance, error)
	LoadCurrentPointsByUserID(ctx context.Context, userID int) (float64, error)
	LoadWithdrawalsByUserID(ctx context.Context, userID int) ([]model.Withdrawals, error)
	SaveWithdrawals(ctx context.Context, userID int, withdrawls model.Withdrawals) error
	SaveAccrualStatusByOrderNumber(ctx context.Context, orderNumber int, status string) error
	SaveAccrualByOrderNumber(ctx context.Context, orderNumber int, accrual float64) error
	SaveNewPoints(ctx context.Context, userID int, newPoints float64) error
	LoadWithdrawn(ctx context.Context, userID int) (float64, error)
	SaveNewWithdrawn(ctx context.Context, userID int, withdrawn float64) error
}

type StorageStrInt struct {
	AllStorageInt
}

func NewStorage(db *sqlx.DB) *StorageStrInt {
	return &StorageStrInt{
		AllStorageInt: postgres.NewAllPostgres(db),
	}
}
