package service

import (
	"context"

	"github.com/golovanevvs/gophermart/internal/customerrors"
	"github.com/golovanevvs/gophermart/internal/model"
)

func (os *orderServiceStr) GetOrders(ctx context.Context, userID int) ([]model.Order, customerrors.CustomError) {
	orders, err := os.st.LoadOrderByUserID(ctx, userID)
	if err != nil {
		return nil, customerrors.New(err, customerrors.DBError500)
	}

	if len(orders) == 0 {
		return nil, customerrors.New(nil, customerrors.EmptyOrder204)
	}

	return orders, customerrors.New(nil, "")
}
