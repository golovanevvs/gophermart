package service

import (
	"context"
	"fmt"

	"github.com/golovanevvs/gophermart/internal/customerrors"
	"github.com/golovanevvs/gophermart/internal/model"
)

func (os *orderServiceStr) GetOrders(ctx context.Context, userID int) ([]model.Order, error) {
	orders, err := os.st.LoadOrderByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%v: %v", customerrors.DBError500, err.Error())
	}

	if len(orders) == 0 {
		return nil, fmt.Errorf("%v", customerrors.EmptyOrder204)
	}

	return orders, nil
}
