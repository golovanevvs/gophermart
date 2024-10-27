package service

import (
	"context"

	"github.com/golovanevvs/gophermart/internal/model"
)

func (os *orderServiceStr) GetOrderNumbers(ctx context.Context, userID int) ([]model.Order, error) {

	return []model.Order{}, nil
}
