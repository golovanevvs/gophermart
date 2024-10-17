package postgres

import (
	"context"

	"github.com/golovanevvs/gophermart/internal/model"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type allPostgresStr struct {
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

func NewAllPostgres(db *sqlx.DB) *allPostgresStr {
	return &allPostgresStr{
		db: db,
	}
}

func (ap *allPostgresStr) CreateUser(ctx context.Context, user model.User) (int, error) {
	var userID int

	row := ap.db.QueryRowContext(ctx, `
	INSERT INTO account
		(login, password_hash)
	VALUES
		($1, $2)
	RETURNING user_id
	`, user.Login, user.PasswordHash)
	if err := row.Scan(&userID); err != nil {
		return 0, err
	}

	return userID, nil
}

func (ap *allPostgresStr) GetUserByLoginPasswordHash(ctx context.Context, login, passwordHash string) (model.User, error) {
	row := ap.db.QueryRowContext(ctx, `
	SELECT user_id FROM %s
	WHERE login=$1 AND password_hash=$2;
	`, login, passwordHash)

	var user model.User

	err := row.Scan(&user.UserID)
	if err != nil {
		return model.User{}, err
	}

	user.Login = login

	return user, nil
}
