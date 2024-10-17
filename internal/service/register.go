package service

import (
	"context"

	"github.com/golovanevvs/gophermart/internal/model"
)

func (as *authServiceStr) CreateUser(ctx context.Context, user model.User) (int, error) {
	user.PasswordHash = genPasswordHash(user.Password)
	return as.st.CreateUser(ctx, user)
}
