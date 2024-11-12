package service

import (
	"context"

	"github.com/golovanevvs/gophermart/internal/model"
)

func (as *authServiceStr) GetBalance(ctx context.Context, userID int) (model.Balance, error) {
	return as.st.LoadBalanceByUserID(ctx, userID)
}
