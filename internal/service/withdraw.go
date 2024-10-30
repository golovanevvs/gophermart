package service

import (
	"context"
	"strconv"

	"github.com/golovanevvs/gophermart/internal/customerrors"
)

func (os *orderServiceStr) Withdraw(ctx context.Context, userID int, withdrawOrderNumber string, sum int) customerrors.CustomError {
	// преобразование номера заказа в число
	withdrawOrderNumberInt, err := strconv.Atoi(withdrawOrderNumber)
	if err != nil {
		return customerrors.New(err, customerrors.InvalidOrderNumberNotInt422)
	}

	// проверка номера заказа алгоритмом Луна
	if err := checkOrderNumberByLuhn(withdrawOrderNumberInt); err != nil {
		return customerrors.New(err, customerrors.InvalidOrderNumber422)
	}

	// загрузка текущего баланса пользователя
	currentPoints, err := os.st.LoadCurrentPointsByUserID(ctx, userID)
	if err != nil {
		return customerrors.New(err, customerrors.DBError500)
	}

	// проверка, что на счету достаточно средств
	if sum > currentPoints {
		return customerrors.New(nil, customerrors.NotEnoughPoints402)
	}

	// TODO: добавить списание баллов со счёта пользователя

	return customerrors.New(nil, "")
}
