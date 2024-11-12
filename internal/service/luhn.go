package service

import (
	"errors"
	"strconv"
)

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
