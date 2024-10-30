package postgres

import (
	"context"
	"time"

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

	// пингование БД
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// удаление таблиц БД
	err = dropTables(db)
	if err != nil {
		return nil, err
	}

	// создание таблиц БД
	err = createTables(db)
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

func createTables(db *sqlx.DB) error {
	// таймаут 5 секунд
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// создание таблицы account, если не существует
	_, err := db.ExecContext(ctx, `
	CREATE TABLE IF NOT EXISTS account (
		user_id SERIAL PRIMARY KEY,
		login VARCHAR(250) UNIQUE NOT NULL,
		password_hash VARCHAR(250) NOT NULL
	);
	`)
	if err != nil {
		return err
	}

	// создание таблицы orders, если не существует
	_, err = db.ExecContext(ctx, `
	CREATE TABLE IF NOT EXISTS orders (
		order_id SERIAL PRIMARY KEY,
		order_number BIGINT UNIQUE,
		order_status VARCHAR(250) NOT NULL,
		uploaded_at TIMESTAMPTZ,
		user_id INT NOT NULL,
		FOREIGN KEY (user_id) REFERENCES account(user_id) ON DELETE CASCADE
	);
	`)
	if err != nil {
		return err
	}

	// создание таблицы balance, если не существует
	_, err = db.ExecContext(ctx, `
	CREATE TABLE IF NOT EXISTS balance (
		current_points INT DEFAULT 0,
		withdrawn INT DEFAULT 0,
		user_id INT,
		FOREIGN KEY (user_id) REFERENCES account(user_id) ON DELETE CASCADE
	);
	`)
	if err != nil {
		return err
	}

	// создание таблицы accrual, если не существует
	_, err = db.ExecContext(ctx, `
	CREATE TABLE IF NOT EXISTS accrual (
		accrual_points INT,
		accrual_at TIMESTAMPTZ,
		order_id INT,
		user_id INT,
		FOREIGN KEY (order_id) REFERENCES orders(order_id) ON DELETE CASCADE,
		FOREIGN KEY (user_id) REFERENCES account(user_id) ON DELETE CASCADE
	);
	`)

	return nil
}

func dropTables(db *sqlx.DB) error {
	// таймаут 5 секунд
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// удаление таблиц БД
	_, err := db.ExecContext(ctx, `
	DROP TABLE IF EXISTS accrual;
	DROP TABLE IF EXISTS balance;
	DROP TABLE IF EXISTS orders;
	DROP TABLE IF EXISTS account;
	`)
	if err != nil {
		return err
	}

	return nil
}

func (ap *allPostgresStr) SaveUser(ctx context.Context, user model.User) (int, error) {

	row := ap.db.QueryRowContext(ctx, `
	INSERT INTO account
	(login, password_hash)
	VALUES
	($1, $2)
	RETURNING user_id;
	`, user.Login, user.PasswordHash)

	var userID int
	if err := row.Scan(&userID); err != nil {
		return 0, err
	}

	return userID, nil
}

func (ap *allPostgresStr) LoadUserByLoginPasswordHash(ctx context.Context, login, passwordHash string) (model.User, error) {
	row := ap.db.QueryRowContext(ctx, `
	SELECT user_id FROM account
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

// func (ap *allPostgresStr) LoadPointsByUserID(ctx context.Context, userID int) (int, error) {
// 	row := ap.db.QueryRowContext(ctx, `
// 	SELECT points FROM account
// 	WHERE user_id=$1;
// 	`, userID)

// 	var points int
// 	err := row.Scan(&points)
// 	if err != nil {
// 		return 0, err
// 	}

// 	return points, nil
// }

func (ap *allPostgresStr) SaveOrderNumberByUserID(ctx context.Context, userID int, orderNumber int) (int, error) {
	row := ap.db.QueryRowContext(ctx, `
	INSERT INTO orders
		(order_number, user_id, order_status, uploaded_at)
	VALUES
		($1, $2, 'NEW', NOW())
	RETURNING order_id;
	`, orderNumber, userID)

	var orderID int
	if err := row.Scan(&orderID); err != nil {
		return 0, err
	}

	return orderID, nil
}

func (ap *allPostgresStr) LoadUserIDByOrderNumber(ctx context.Context, orderNumber int) (int, error) {
	row := ap.db.QueryRowContext(ctx, `
	SELECT user_id FROM orders
	WHERE order_number=$1;
	`, orderNumber)

	var userID int
	err := row.Scan(&userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (ap *allPostgresStr) LoadOrderByUserID(ctx context.Context, userID int) ([]model.Order, error) {
	orders := make([]model.Order, 0)

	rows, err := ap.db.QueryContext(ctx, `
	SELECT
	order_id, order_number, order_status, uploaded_at
	FROM orders
	WHERE user_id = $1
	ORDER BY uploaded_at DESC;
	`, userID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var order model.Order
		err = rows.Scan(&order.OrderID, &order.OrderNumber, &order.OrderStatus, &order.UploadedAt)
		if err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (ap *allPostgresStr) LoadBalanceByUserID(ctx context.Context, userID int) (model.Balance, error) {
	row := ap.db.QueryRowContext(ctx, `
	SELECT current_points, withdrawn FROM balance
	WHERE user_id=$1;
	`, userID)

	var currentPoints, withdrawn int

	err := row.Scan(&currentPoints, &withdrawn)
	if err != nil {
		return model.Balance{}, err
	}

	balance := model.Balance{
		CurrentPoints: currentPoints,
		Withdrawn:     withdrawn,
	}

	return balance, nil
}

func (ap *allPostgresStr) LoadCurrentPointsByUserID(ctx context.Context, userID int) (int, error) {
	row := ap.db.QueryRowContext(ctx, `
	SELECT current_balance FROM balance
	WHERE user_id=$1;
	`, userID)

	var currentPoints int

	err := row.Scan(&currentPoints)
	if err != nil {
		return 0, err
	}

	return currentPoints, nil
}
