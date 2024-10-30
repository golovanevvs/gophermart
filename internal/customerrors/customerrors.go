package customerrors

import (
	"errors"
	"fmt"
)

type customErr string

const (
	OrderAlredyUploadedThisUser200  customErr = "номер заказа уже был загружен этим пользователем"
	EmptyOrder204                   customErr = "нет данных для ответа"
	InvalidRequest400               customErr = "неверный формат запроса"
	InvalidContentType400           customErr = "неверный Content-Type"
	DBInvalidLoginPassword401       customErr = "неверная пара логин/пароль"
	JWTWrongSingingMethod401        customErr = "неверный метод подписи"
	JWTParseError401                customErr = "ошибка при чтении JWT"
	JWTInvalidToken401              customErr = "невалидный токен"
	NotEnoughPoints402              customErr = "на счету недостаточно средств"
	OrderAlredyUploadedOtherUser409 customErr = "номер заказа уже был загружен другим пользователем"
	DBBusyLogin409                  customErr = "ошибка БД: логин уже занят"
	InvalidOrderNumber422           customErr = "Неверный формат номера заказа"
	InvalidOrderNumberNotInt422     customErr = "Неверный формат номера заказа: не соответствует типу int"
	DBError500                      customErr = "ошибка БД"
	InternalServerError500          customErr = "внутренняя ошибка сервера"
	DecodeJSONError500              customErr = "ошибка десериализации JSON"
)

type CustomError struct {
	IsError   bool
	Err       error
	CustomErr customErr
	AllErr    error
}

func New(err error, customErr customErr) CustomError {
	var customError CustomError

	if err != nil || len(customErr) > 0 {
		customError.IsError = true
	}

	customError.Err = err
	customError.CustomErr = customErr
	customError.AllErr = errors.New(fmt.Sprintf("ошибка: %v; %v", customErr, err))

	return customError
}
