package service

import (
	"context"
	"strings"

	"github.com/golovanevvs/gophermart/internal/customerrors"
)

func (os *orderServiceStr) UploadOrder(ctx context.Context, userID int, orderNumber int) (int, customerrors.CustomError) {
	// проверка номера заказа алгоритмом Луна
	if err := checkOrderNumberByLuhn(orderNumber); err != nil {
		return 0, customerrors.New(err, customerrors.InvalidOrderNumber422)
	}

	orderID, err := os.st.SaveOrderNumberByUserID(ctx, userID, orderNumber)
	if err != nil {
		// если номер заказа уже есть в БД
		if strings.Contains(err.Error(), " 23505") {
			// получение userID (newUserID), которому принадлежит номер заказа
			newUserID, err := os.st.LoadUserIDByOrderNumber(ctx, orderNumber)
			// если возникла непредвиденная ошибка
			if err != nil {
				return 0, customerrors.New(err, customerrors.DBError500)
			}
			if newUserID == userID {
				// если номер заказа принадлежит текущему пользователю
				return 0, customerrors.New(err, customerrors.OrderAlredyUploadedThisUser200)
			} else {
				// если номер заказа принадлежит другому пользователю
				return 0, customerrors.New(err, customerrors.OrderAlredyUploadedOtherUser409)
			}
		}
		// если возникла другая ошибка
		return 0, customerrors.New(err, customerrors.DBError500)
	}

	return orderID, customerrors.New(nil, "")
}
