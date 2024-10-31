package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/golovanevvs/gophermart/internal/customerrors"
	"github.com/golovanevvs/gophermart/internal/model"
)

func (as *authServiceStr) CreateUser(ctx context.Context, user model.User) (int, error) {
	// хеширование пароля
	user.PasswordHash = genPasswordHash(user.Password)

	// запуск функции БД: сохранение нового пользователя
	userID, err := as.st.SaveUser(ctx, user)
	if err != nil {
		if strings.Contains(err.Error(), " 23505") {
			return 0, fmt.Errorf("%v: %v", customerrors.DBBusyLogin409, err.Error())
		}
		return 0, fmt.Errorf("%v: %v", customerrors.DBError500, err.Error())
	}

	return userID, nil
}
