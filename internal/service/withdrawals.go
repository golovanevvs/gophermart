package service

import (
	"context"
	"fmt"

	"github.com/golovanevvs/gophermart/internal/customerrors"
	"github.com/golovanevvs/gophermart/internal/model"
)

func (os *orderServiceStr) Withdrawals(ctx context.Context, userID int) ([]model.Withdrawals, error) {
	// обращение к БД
	withdrawals, err := os.st.LoadWithdrawalsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%v: %v", customerrors.DBError500, err.Error())
	}

	// если нет ни одного списания
	if len(withdrawals) == 0 {
		return nil, fmt.Errorf("%v", customerrors.EmptyWithdrawals204)
	}

	return withdrawals, nil
}
