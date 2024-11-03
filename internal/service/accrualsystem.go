package service

import (
	"context"

	"github.com/golovanevvs/gophermart/internal/model"
)

func (as *accrualSystemServiceStr) GetOrderFromAPIByOrderNumber(ctx context.Context, orderNumber int) (model.AccrualSystem, error) {
	accrualSystem, err := as.as.GetAPIOrders(ctx, orderNumber)
	// TODO: доработать сервис
	if err != nil {
		return model.AccrualSystem{}, err
	}

	return accrualSystem, nil
}
