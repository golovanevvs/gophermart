package customerrors

import "errors"

type customErr string

const (
	Err400                 customErr = "неверный формат запроса"
	Err401                 customErr = "неверная пара логин/пароль"
	BusyLogin409           customErr = "логин уже занят"
	InternalServerError500 customErr = "внутренняя ошибка сервера"
)

type CustomError struct {
	IsError   bool
	Err       error
	CustomErr error
}

func New(err error, customErr customErr) CustomError {
	var customError CustomError

	customError.IsError = false

	if err != nil || len(customErr) > 0 {
		customError.IsError = true
	}

	customError.Err = err
	customError.CustomErr = errors.New(string(customErr))

	return customError
}
