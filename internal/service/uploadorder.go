package service

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/golovanevvs/gophermart/internal/customerrors"
)

func (os *orderServiceStr) UploadOrder(ctx context.Context, userID int, orderNumber int) (int, customerrors.CustomError) {
	// проверка номера заказа алгоритмом Луна
	if err := checkOrderNumberByLuhn(orderNumber); err != nil {
		return 0, customerrors.New(err, customerrors.LuhnInvalid422)
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

func checkOrderNumberByLuhn(orderNumber int) error {
	// преобразование числа в строку
	orderNumberString := strconv.Itoa(orderNumber)

	// проверка, что номер заказа содержит больше одной цифры
	if len(orderNumberString) <= 1 {
		return errors.New("минимальная длина должна быть больше одной цифры")
	}

	// проверка корректности последовательности цифр с использованием алгоритма Луна
	sum := 0
	isSecond := false
	for i := len(orderNumberString) - 1; i >= 0; i-- {
		digit, _ := strconv.Atoi(string(orderNumberString[i]))

		if isSecond {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
		isSecond = !isSecond
	}
	if sum%10 == 0 {
		return nil
	} else {
		return errors.New("номер заказа не прошёл проверку с использованием алгоритма Луна")
	}
}
