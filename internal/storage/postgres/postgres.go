package postgres

import (
	"context"
	"strconv"
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
		password_hash VARCHAR(250) NOT NULL,
		current_points DOUBLE PRECISION DEFAULT 0,
		withdrawn DOUBLE PRECISION DEFAULT 0
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
		uploaded_at TIMESTAMP WITH TIME ZONE,
		accrual DOUBLE PRECISION DEFAULT 0,
		user_id INT NOT NULL,
		FOREIGN KEY (user_id) REFERENCES account(user_id) ON DELETE CASCADE
	);
	`)
	if err != nil {
		return err
	}

	// создание таблицы withdrawals, если не существует
	_, err = db.ExecContext(ctx, `
	CREATE TABLE IF NOT EXISTS withdrawals (
		withdrawals_id SERIAL PRIMARY KEY,
		new_order VARCHAR(250),
		sum DOUBLE PRECISION,
		processed_at TIMESTAMP WITH TIME ZONE,
		user_id INT,
		FOREIGN KEY (user_id) REFERENCES account(user_id) ON DELETE CASCADE
	);
	`)
	if err != nil {
		return err
	}

	return nil
}

func dropTables(db *sqlx.DB) error {
	// таймаут 5 секунд
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// удаление таблиц БД
	_, err := db.ExecContext(ctx, `
	DROP TABLE IF EXISTS withdrawals;
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

	err := row.Scan(&user.ID)
	if err != nil {
		return model.User{}, err
	}

	user.Login = login

	return user, nil
}

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
	order_id, order_number, order_status, accrual, uploaded_at
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
		err = rows.Scan(&order.ID, &order.Number, &order.Status, &order.Accrual, &order.UploadedAt)
		if err != nil {
			return nil, err
		}

		order.NumberString = strconv.Itoa(order.Number)

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
	SELECT current_points, withdrawn FROM account
	WHERE user_id=$1;
	`, userID)

	var currentPoints, withdrawn float64

	err := row.Scan(&currentPoints, &withdrawn)
	if err != nil {
		return model.Balance{}, err
	}

	balance := model.Balance{
		Current:   currentPoints,
		Withdrawn: withdrawn,
	}

	return balance, nil
}

func (ap *allPostgresStr) LoadCurrentPointsByUserID(ctx context.Context, userID int) (float64, error) {
	row := ap.db.QueryRowContext(ctx, `
	SELECT current_points FROM account
	WHERE user_id=$1;
	`, userID)

	var currentPoints float64

	err := row.Scan(&currentPoints)
	if err != nil {
		return 0, err
	}

	return currentPoints, nil
}

func (ap *allPostgresStr) LoadWithdrawalsByUserID(ctx context.Context, userID int) ([]model.Withdrawals, error) {
	withdrawals := make([]model.Withdrawals, 0)

	rows, err := ap.db.QueryContext(ctx, `
	SELECT new_order, sum, processed_at FROM withdrawals
	WHERE user_id=$1;
	`, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var withdrawn model.Withdrawals
		err := rows.Scan(&withdrawn.NewOrderNumber, &withdrawn.Sum, &withdrawn.ProcessedAt)
		if err != nil {
			return nil, err
		}

		withdrawals = append(withdrawals, withdrawn)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return withdrawals, nil
}

func (ap *allPostgresStr) SaveWithdrawals(ctx context.Context, userID int, withdrawls model.Withdrawals) error {
	_, err := ap.db.ExecContext(ctx, `
	INSERT INTO withdrawals
		(new_order, sum, processed_at, user_id)
	VALUES
		($1, $2, $3, $4);
	`, withdrawls.NewOrderNumber, withdrawls.Sum, withdrawls.ProcessedAt, userID)
	if err != nil {
		return err
	}

	return nil
}

func (ap *allPostgresStr) SaveAccrualStatusByOrderNumber(ctx context.Context, orderNumber int, status string) error {
	_, err := ap.db.ExecContext(ctx, `
	UPDATE orders
	SET	order_status = $1
	WHERE order_number = $2;
	`, status, orderNumber)
	if err != nil {
		return err
	}
	return nil
}

func (ap *allPostgresStr) SaveAccrualByOrderNumber(ctx context.Context, orderNumber int, accrual float64) error {
	_, err := ap.db.ExecContext(ctx, `
	UPDATE orders
	SET accrual = $1
	WHERE order_number = $2;
	`, accrual, orderNumber)
	if err != nil {
		return err
	}
	return nil
}

func (ap *allPostgresStr) SaveNewPoints(ctx context.Context, userID int, newPoints float64) error {
	_, err := ap.db.ExecContext(ctx, `
	UPDATE account
	SET current_points = $1
	WHERE user_id = $2;
	`, newPoints, userID)
	if err != nil {
		return err
	}
	return nil
}

func (ap *allPostgresStr) LoadWithdrawn(ctx context.Context, userID int) (float64, error) {
	row := ap.db.QueryRowContext(ctx, `
	SELECT withdrawn FROM account
	WHERE user_id = $1;
	`, userID)

	var withdrawn float64

	err := row.Scan(&withdrawn)
	if err != nil {
		return 0, err
	}

	return withdrawn, nil
}

func (ap *allPostgresStr) SaveNewWithdrawn(ctx context.Context, userID int, withdrawn float64) error {
	_, err := ap.db.ExecContext(ctx, `
	UPDATE account
	SET withdrawn = $1
	WHERE user_id = $2;
	`, withdrawn, userID)
	if err != nil {
		return err
	}
	return nil
}
