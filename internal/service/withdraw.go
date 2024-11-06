package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/golovanevvs/gophermart/internal/customerrors"
	"github.com/golovanevvs/gophermart/internal/model"
)

func (os *orderServiceStr) Withdraw(ctx context.Context, userID int, withdrawOrderNumber string, sum float64) error {
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

	// загрузка withdrawn
	currentWithdrawn, err := os.st.LoadWithdrawn(ctx, userID)
	if err != nil {
		return fmt.Errorf("%v: %v", customerrors.DBError500, err.Error())
	}

	// обновление баланса
	newPoints := currentPoints - sum
	err = os.st.SaveNewPoints(ctx, userID, newPoints)
	if err != nil {
		return fmt.Errorf("%v: %v", customerrors.DBError500, err.Error())
	}

	newWithdrawn := currentWithdrawn + sum
	err = os.st.SaveNewWithdrawn(ctx, userID, newWithdrawn)
	if err != nil {
		return fmt.Errorf("%v: %v", customerrors.DBError500, err.Error())
	}

	// обновление withdrawals
	newWithdrawals := model.Withdrawals{
		NewOrderNumber: withdrawOrderNumber,
		Sum:            sum,
		ProcessedAt:    time.Now(),
	}
	err = os.st.SaveWithdrawals(ctx, newWithdrawals)
	if err != nil {
		return fmt.Errorf("%v: %v", customerrors.DBError500, err.Error())
	}

	return nil
}
