package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/golovanevvs/gophermart/internal/customerrors"
)

func (os *orderServiceStr) UploadOrder(ctx context.Context, userID int, orderNumber int) (int, error) {
	// проверка номера заказа алгоритмом Луна
	if err := checkOrderNumberByLuhn(orderNumber); err != nil {
		return 0, fmt.Errorf("%v: %v", customerrors.InvalidOrderNumber422, err.Error())
	}

	// сохранение в БД, обработка ошибок
	orderID, err := os.st.SaveOrderNumberByUserID(ctx, userID, orderNumber)
	if err != nil {
		// если номер заказа уже есть в БД
		if strings.Contains(err.Error(), " 23505") {
			// получение userID (newUserID), которому принадлежит номер заказа
			newUserID, err := os.st.LoadUserIDByOrderNumber(ctx, orderNumber)
			// если возникла непредвиденная ошибка
			if err != nil {
				return 0, fmt.Errorf("%v: %v", customerrors.DBError500, err.Error())
			}
			if newUserID == userID {
				// если номер заказа принадлежит текущему пользователю
				return 0, fmt.Errorf("%v", customerrors.OrderAlredyUploadedThisUser200)
			} else {
				// если номер заказа принадлежит другому пользователю
				return 0, fmt.Errorf("%v", customerrors.OrderAlredyUploadedOtherUser409)
			}
		}
		// если возникла другая ошибка
		return 0, fmt.Errorf("%v: %v", customerrors.DBError500, err.Error())
	}

	go func(userID int, orderNumber int) {
		err = os.as.GetOrderFromAS(userID, orderNumber)
		if err != nil {
			//TODO: добавить логгирование
			fmt.Printf("Ошибка в AS: %v", err.Error())
		}
	}(userID, orderNumber)

	return orderID, nil
}
