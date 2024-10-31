package service

import (
	"context"
	"fmt"
	"strconv"

	"github.com/golovanevvs/gophermart/internal/customerrors"
)

func (os *orderServiceStr) Withdraw(ctx context.Context, userID int, withdrawOrderNumber string, sum int) error {
	// преобразование номера заказа в число
	withdrawOrderNumberInt, err := strconv.Atoi(withdrawOrderNumber)
	if err != nil {
		return fmt.Errorf("%v: %v", customerrors.InvalidOrderNumberNotInt422, err.Error())
	}

	// проверка номера заказа алгоритмом Луна
	if err := checkOrderNumberByLuhn(withdrawOrderNumberInt); err != nil {
		return fmt.Errorf("%v: %v", customerrors.InvalidOrderNumber422, err.Error())
	}

	// загрузка текущего баланса пользователя
	currentPoints, err := os.st.LoadCurrentPointsByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("%v: %v", customerrors.DBError500, err.Error())
	}

	// проверка, что на счету достаточно средств
	if sum > currentPoints {
		return fmt.Errorf("%v", customerrors.NotEnoughPoints402)
	}

	// TODO: добавить списание баллов со счёта пользователя

	return nil
}
