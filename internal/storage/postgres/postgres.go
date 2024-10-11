package postgres

import (
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type AuthPostgresStr struct {
	db *sqlx.DB
}

type RegisterPostgresStr struct {
	db *sqlx.DB
}

func NewPostgres(databaseURI string) (*sqlx.DB, error) {
	db, err := sqlx.Open("pgx", databaseURI)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgresStr {
	return &AuthPostgresStr{
		db: db,
	}
}

func NewRegisterPostgres(db *sqlx.DB) *RegisterPostgresStr {
	return &RegisterPostgresStr{
		db: db,
	}
}

func (ap *AuthPostgresStr) CreateUser() string {
	return fmt.Sprintf("CreateUser")
}
