package service

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/golovanevvs/gophermart/internal/customerrors"
	"github.com/golovanevvs/gophermart/internal/model"
)

func (as *accrualSystemServiceStr) ProcAccrual(userID int, orderNumber int) {
	ctx := context.Background()
	accrual, err := as.GetAccrual(ctx, userID, orderNumber)
	if err != nil {
		fmt.Printf("%v: %v", customerrors.ASError, err.Error())
		return
	}
	if accrual > 0 {
		err := as.UpdateBalance(ctx, userID, accrual)
		if err != nil {
			fmt.Printf("%v: %v", customerrors.ASError, err.Error())
			return
		}
	}
}

func (as *accrualSystemServiceStr) GetAccrual(ctx context.Context, userID int, orderNumber int) (float64, error) {
	var accrualSystem model.AccrualSystem
	var err error
	var interval time.Duration

	for {
		interval = time.Second * 5
		accrualSystem, err = as.as.GetAPIOrders(ctx, orderNumber)
		if err != nil {
			switch {
			case strings.Contains(err.Error(), customerrors.InvalidContentType400):
				// статус заказа не менять
				// отобразить ошибку в логах
				// прервать
				return 0, err
			case strings.Contains(err.Error(), customerrors.DecodeJSONError500):
				// статус заказа не менять
				// отобразить ошибку в логах
				// прервать
				fmt.Printf("%v: %v", customerrors.ASError, err.Error())
				return 0, err
			case strings.Contains(err.Error(), customerrors.ASOrderNotRegistered204):
				// статус заказа не менять
				// отобразить ошибку в логах
				// прервать
				fmt.Printf("%v: %v", customerrors.ASError, err.Error())
				return 0, err
			case strings.Contains(err.Error(), customerrors.AtoiError500):
				// статус заказа не менять
				// отобразить ошибку в логах
				// прервать
				fmt.Printf("%v: %v", customerrors.ASError, err.Error())
				return 0, err
			case strings.Contains(err.Error(), customerrors.ResponseBodyError500):
				// статус заказа не менять
				// отобразить ошибку в логах
				// прервать
				fmt.Printf("%v: %v", customerrors.ASError, err.Error())
				return 0, err
			case strings.Contains(err.Error(), customerrors.ASTooManyRequests429):
				// статус заказа не менять
				// отобразить ошибку в логах
				// повторить серию запросов через время Retry-After
				fmt.Printf("%v: %v", customerrors.ASError, err.Error())
				interval = time.Second * time.Duration(accrualSystem.RetryAfter)
			case strings.Contains(err.Error(), customerrors.InternalServerError500):
				// статус заказа не менять
				// отобразить ошибку в логах
				// прервать
				fmt.Printf("%v: %v", customerrors.ASError, err.Error())
				return 0, err
			}
		} else {
			switch accrualSystem.Status {
			case "REGISTERED":
				// статус заказа не менять (остаётся NEW)
			case "INVALID":
				// поменять статус заказа - INVALID
				// прервать
				as.st.SaveAccrualStatusByOrderNumber(ctx, orderNumber, "INVALID")
				return 0, nil
			case "PROCESSING":
				// поменять статус заказа - PROCESSING
				as.st.SaveAccrualStatusByOrderNumber(ctx, orderNumber, "PROCESSING")
			case "PROCESSED":
				// поменять статус заказа - PROCESSED
				// записать accrual
				// обновить баланс
				// прервать
				as.st.SaveAccrualStatusByOrderNumber(ctx, orderNumber, "PROCESSED")
				as.st.SaveAccrualByOrderNumber(ctx, orderNumber, accrualSystem.Accrual)
				return accrualSystem.Accrual, nil
			}
		}
		time.Sleep(interval)
	}
}

func (as *accrualSystemServiceStr) UpdateBalance(ctx context.Context, userID int, accrual float64) error {
	currentPoints, err := as.st.LoadCurrentPointsByUserID(ctx, userID)
	if err != nil {
		return err
	}
	accrualRound := math.Round(accrual*100) / 100
	newPoints := currentPoints + accrualRound
	err = as.st.SaveNewPoints(ctx, userID, newPoints)
	if err != nil {
		return err
	}
	return nil
}
