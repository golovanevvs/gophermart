package service

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const (
	TOKEN_EXP  = time.Hour * 3
	SECRET_KEY = "sskey"
)

type claims struct {
	jwt.RegisteredClaims
	UserID int
}

func (as *authServiceStr) GenToken(ctx context.Context, login, password string) (string, error) {
	user, err := as.st.GetUserByLoginPasswordHash(ctx, login, genPasswordHash(password))
	if err != nil {
		return "", nil
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TOKEN_EXP)),
		},
		UserID: user.UserID,
	})
	tokenString, err := token.SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (as *authServiceStr) ParseToken(tokenString string) (int, error) {
	claims := &claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("неверный метод подписи")
		}
		return []byte(SECRET_KEY), nil
	})
	if err != nil {
		return 0, err
	}

	if !token.Valid {
		return 0, errors.New("невалидный токен")
	}

	return claims.UserID, nil
}
