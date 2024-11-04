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
	var err error

	ctx := context.Background()

	for {
		fmt.Printf("accrualSystem.Status: %v\n", accrualSystem.Status)
		if accrualSystem.Status == "PROCESSED" || accrualSystem.Status == "INVALID" {
			break
		}

		// интервал обращения к сервису
		interval := time.Second * 5

		// для отслеживания измнения статуса
		newStatus := ""

		// запуск обращения к сервису
		accrualSystem, err = as.as.GetAPIOrders(ctx, orderNumber)
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
			fmt.Printf("Сохранение статуса в БД: %v...\n", accrualSystem.Status)
			as.st.SaveAccrualStatusByOrderNumber(ctx, orderNumber, accrualSystem.Status)
			fmt.Printf("Сохранение статуса в БД завершено: %v\n", accrualSystem.Status)
			newStatus = accrualSystem.Status
		}

		// задержка перед повторным обращением
		time.Sleep(interval)
	}

	// сохранение начисления баллов или статуса INVALID в БД
	if accrualSystem.Status == "PROCESSED" {
		fmt.Printf("Сохранение accrualSystem в БД: %v...\n", accrualSystem)
		as.st.SaveAccrualByOrderNumber(ctx, accrualSystem)
		fmt.Printf("Сохранение accrualSystem в БД завершено: %v\n", accrualSystem)
		if accrualSystem.Accrual > 0 {
			fmt.Printf("Получение текущего баланса из БД userID: %v...\n", userID)
			currentPoints, err := as.st.LoadCurrentPointsByUserID(ctx, userID)
			fmt.Printf("Получение текущего баланса из БД userID завершено: %v, баланс: %v\n", userID, currentPoints)
			if err != nil {
				return err
			}
			newPoints := currentPoints + accrualSystem.Accrual
			fmt.Printf("Сохранение нового баланса в БД: %v...\n", newPoints)
			err = as.st.SaveNewPoints(ctx, userID, newPoints)
			fmt.Printf("Сохранение нового баланса в БД завершено: %v\n", newPoints)
			if err != nil {
				fmt.Printf("Ошибка при сохранении нового баланса в БД: %v", err.Error())
				return err
			}
		}
	} else {
		fmt.Printf("Сохранение статуса INVALID в БД: %v...\n", accrualSystem.Status)
		as.st.SaveAccrualStatusByOrderNumber(ctx, orderNumber, accrualSystem.Status)
		fmt.Printf("Сохранение статуса INVALID в БД завершено: %v\n", accrualSystem.Status)
	}

	return nil
}
