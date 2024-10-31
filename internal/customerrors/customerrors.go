package customerrors

const (
	OrderAlredyUploadedThisUser200  = "номер заказа уже был загружен этим пользователем"
	EmptyOrder204                   = "нет данных для ответа"
	InvalidRequest400               = "неверный формат запроса"
	InvalidContentType400           = "неверный Content-Type"
	DBInvalidLoginPassword401       = "неверная пара логин/пароль"
	JWTWrongSingingMethod401        = "неверный метод подписи"
	JWTParseError401                = "ошибка при чтении JWT"
	JWTInvalidToken401              = "невалидный токен"
	NotEnoughPoints402              = "на счету недостаточно средств"
	OrderAlredyUploadedOtherUser409 = "номер заказа уже был загружен другим пользователем"
	DBBusyLogin409                  = "ошибка БД: логин уже занят"
	InvalidOrderNumber422           = "Неверный формат номера заказа"
	InvalidOrderNumberNotInt422     = "Неверный формат номера заказа: не соответствует типу int"
	DBError500                      = "ошибка БД"
	InternalServerError500          = "внутренняя ошибка сервера"
	DecodeJSONError500              = "ошибка десериализации JSON"
)
