package customerrors

import "errors"

type customErr string

const (
	Err400                 customErr = "неверный формат запроса"
	Err401                 customErr = "неверная пара логин/пароль"
	Err409                 customErr = "логин уже занят"
	InternalServerError500 customErr = "внутренняя ошибка сервера"
)

type CustomError struct {
	Err       error
	CustomErr error
}

func New(err error, customErr customErr) CustomError {
	return CustomError{
		Err:       err,
		CustomErr: errors.New(string(customErr)),
	}
}
