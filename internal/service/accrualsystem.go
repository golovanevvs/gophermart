package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/golovanevvs/gophermart/internal/customerrors"
	"github.com/golovanevvs/gophermart/internal/model"
)

// TODO: доработать сервис
func (as *accrualSystemServiceStr) GetOrderFromASByOrderNumber(ctx context.Context, orderNumber int) (accrualSystem model.AccrualSystem, err error) {
	for accrualSystem.Status != "INVALID" || accrualSystem.Status != "PROCESSED" {
		// интервал обращения к системе
		interval := time.Second * 5

		accrualSystem, err = as.as.GetAPIOrders(ctx, orderNumber)
		if err != nil {
			// если превышено количество запросов к сервису
			if strings.Contains(err.Error(), customerrors.ASTooManyRequests429) {
				interval = time.Duration(accrualSystem.RetryAfter)
			} else {
				err = fmt.Errorf("ошибка в сервисе по взаимодействию с системой расчёта начислений баллов: %v", err.Error())

				return
			}
		}

		time.Sleep(interval)
	}

	return
}
