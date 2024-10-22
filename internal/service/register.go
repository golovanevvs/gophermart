package service

import (
	"context"
	"strings"

	"github.com/golovanevvs/gophermart/internal/customerrors"
	"github.com/golovanevvs/gophermart/internal/model"
)

func (as *authServiceStr) CreateUser(ctx context.Context, user model.User) (int, customerrors.CustomError) {
	// хеширование пароля
	user.PasswordHash = genPasswordHash(user.Password)

	// запуск функции БД: сохранение нового пользователя
	userID, err := as.st.SaveUser(ctx, user)
	if err != nil {
		if strings.Contains(err.Error(), "Unique") {
			return 0, customerrors.New(err, customerrors.DBBusyLogin409)
		}
		return 0, customerrors.New(err, customerrors.DBError500)
	}

	return userID, customerrors.New(nil, "")
}
