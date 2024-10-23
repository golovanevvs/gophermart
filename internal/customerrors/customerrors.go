package customerrors

import (
	"errors"
	"fmt"
)

type customErr string

const (
	Err400                    customErr = "неверный формат запроса"
	DBInvalidLoginPassword401 customErr = "неверная пара логин/пароль"
	DBError500                customErr = "ошибка БД"
	DBBusyLogin409            customErr = "ошибка БД: логин уже занят"
	InternalServerError500    customErr = "внутренняя ошибка сервера"
	InvalidContentType400     customErr = "неверный Content-Type"
	DecodeJSONError500        customErr = "ошибка десериализации JSON"
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
