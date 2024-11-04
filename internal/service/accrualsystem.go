package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/golovanevvs/gophermart/internal/customerrors"
	"github.com/golovanevvs/gophermart/internal/model"
)

func (as *accrualSystemServiceStr) GetOrderFromAS(userID int, orderNumber int) error {
	var accrualSystem model.AccrualSystem

	ctx := context.Background()

	for accrualSystem.Status != "INVALID" || accrualSystem.Status != "PROCESSED" {
		// интервал обращения к сервису
		interval := time.Second * 5

		// для отслеживания измнения статуса
		newStatus := ""

		// запуск обращения к сервису
		accrualSystem, err := as.as.GetAPIOrders(ctx, orderNumber)
		if err != nil {
			switch {

			// если заказ не зарегистрирован в системе
			case strings.Contains(err.Error(), customerrors.ASOrderNotRegistered204):
				return fmt.Errorf("%v: %v", customerrors.ASError, err.Error())

			// если превышено количество запросов к сервису
			case strings.Contains(err.Error(), customerrors.ASTooManyRequests429):
				interval = time.Duration(accrualSystem.RetryAfter)

			// если произошла внутренняя ошибка сервера сервиса, то продолжаем обращение к сервису
			default:
			}
		}

		// сохранение статуса в БД
		if accrualSystem.Status != newStatus && accrualSystem.Status != "REGISTERED" {
			as.st.SaveAccrualStatusByOrderNumber(ctx, orderNumber, accrualSystem.Status)
			newStatus = accrualSystem.Status
		}

		// задержка перед повторным обращением
		time.Sleep(interval)
	}

	// сохранение начисления баллов или статуса INVALID в БД
	if accrualSystem.Status == "PROCESSED" {
		as.st.SaveAccrualByOrderNumber(ctx, accrualSystem)
		if accrualSystem.Accrual > 0 {
			currentPoints, err := as.st.LoadCurrentPointsByUserID(ctx, userID)
			if err != nil {
				return err
			}
			newPoints := currentPoints + accrualSystem.Accrual
			err = as.st.SaveNewPoints(ctx, userID, newPoints)
			if err != nil {
				return err
			}
		}
	} else {
		as.st.SaveAccrualStatusByOrderNumber(ctx, orderNumber, accrualSystem.Status)
	}

	return nil
}
