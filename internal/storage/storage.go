package storage

import (
	"context"

	"github.com/golovanevvs/gophermart/internal/model"
	"github.com/golovanevvs/gophermart/internal/storage/postgres"
	"github.com/jmoiron/sqlx"
)

type AllStorageInt interface {
	CreateUser(ctx context.Context, user model.User) error
	GetUser(ctx context.Context, login, password string) (model.User, error)
}

type StorageStrInt struct {
	AllStorageInt
}

func NewStorage(db *sqlx.DB) *StorageStrInt {
	return &StorageStrInt{
		AllStorageInt: postgres.NewAllPostgres(db),
	}
}
