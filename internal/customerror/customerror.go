package customerror

import "errors"

type CustomError error

var (
	Err400 CustomError = errors.New("неверный формат запроса")
	Err401 CustomError = errors.New("неверная пара логин/пароль")
	Err409 CustomError = errors.New("логин уже занят")
	Err500 CustomError = errors.New("внутренняя ошибка сервера")
)
