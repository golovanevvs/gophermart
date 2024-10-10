package storage

import (
	"github.com/golovanevvs/gophermart/internal/storage/postgres"
	"github.com/jmoiron/sqlx"
)

type AuthInt interface {
	CreateUser() string
}

type StorageStr struct {
	AuthInt
}

func NewStorage(db *sqlx.DB) *StorageStr {
	return &StorageStr{
		AuthInt: postgres.NewAuthPostgres(db),
	}
}
