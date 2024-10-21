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
	LoadPointsByUserID(ctx context.Context, userID int) (int, error)
	SaveOrderNumberByUserID(ctx context.Context, userID int, orderNumber int) (int, error)
}

type StorageStrInt struct {
	AllStorageInt
}

func NewStorage(db *sqlx.DB) *StorageStrInt {
	return &StorageStrInt{
		AllStorageInt: postgres.NewAllPostgres(db),
	}
}
