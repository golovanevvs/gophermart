package service

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/golovanevvs/gophermart/internal/customerrors"
)

const (
	TOKEN_EXP  = time.Hour * 3
	SECRET_KEY = "sskey"
)

type claims struct {
	jwt.RegisteredClaims
	UserID int
}

func (as *authServiceStr) BuildJWTString(ctx context.Context, login, password string) (string, customerrors.CustomError) {
	user, err := as.st.GetUserByLoginPasswordHash(ctx, login, genPasswordHash(password))
	if err != nil {
		return "", customerrors.New(err, customerrors.InternalServerError500)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TOKEN_EXP)),
		},
		UserID: user.UserID,
	})
	tokenString, err := token.SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", customerrors.New(err, customerrors.Err400)
	}
	return tokenString, customerrors.New(nil)
}

func (as *authServiceStr) GetUserIDFromJWT(tokenString string) (int, error) {
	claims := &claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("неверный метод подписи")
		}
		return []byte(SECRET_KEY), nil
	})
	if err != nil {
		return -1, err
	}

	if !token.Valid {
		return -1, errors.New("невалидный токен")
	}

	return claims.UserID, nil
}
