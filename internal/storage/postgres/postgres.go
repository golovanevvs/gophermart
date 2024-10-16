package postgres

import (
	"context"
	"fmt"

	"github.com/golovanevvs/gophermart/internal/model"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

const (
	usersTable = "account"
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

func (ap *allPostgresStr) CreateUser(ctx context.Context, user model.User) error {
	query := fmt.Sprintf(`
	INSERT INTO %s
		(login, password_hash)
		VALUES
		($1, $2)
	`, usersTable)

	_, err := ap.db.ExecContext(ctx, query, user.Login, user.Password)
	if err != nil {
		return err
	}

	return nil
}

func (ap *allPostgresStr) GetUserByLoginPassword(ctx context.Context, login, password string) (model.User, error) {
	query := fmt.Sprintf(`
	SELECT user_id FROM %s
	WHERE login=$1 AND password_hash=$2;
	`, usersTable)

	row := ap.db.QueryRowContext(ctx, query, login, password)

	var user model.User

	err := row.Scan(&user.UserID)
	if err != nil {
		return model.User{}, err
	}

	user.Login = login

	return user, nil
}
