package service

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/golovanevvs/gophermart/internal/customerrors"
)

const (
	TOKEN_EXP  = time.Hour * 3 // время жизни токена
	SECRET_KEY = "sskey"       // секретный ключ токена
)

// структура утверждений
type claims struct {
	jwt.RegisteredClaims
	UserID int
}

// BuildJWTString создаёт токен и возвращает его в виде строки
func (as *authServiceStr) BuildJWTString(ctx context.Context, login, password string) (string, customerrors.CustomError) {
	// получение пользователя из БД
	user, err := as.st.GetUserByLoginPasswordHash(ctx, login, genPasswordHash(password))
	if err != nil {
		return "", customerrors.New(err, customerrors.InternalServerError500)
	}

	// создание токена
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TOKEN_EXP)),
		},
		UserID: user.UserID,
	})

	// создание строки токена
	tokenString, err := token.SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", customerrors.New(err, customerrors.InternalServerError500)
	}
	return tokenString, customerrors.New(nil, "")
}

// GetUserIDFromJWT возвращает userID из JWT
func (as *authServiceStr) GetUserIDFromJWT(tokenString string) (int, error) {
	claims := &claims{}

	// преобразование строки в токен
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("неверный метод подписи")
		}
		return []byte(SECRET_KEY), nil
	})
	if err != nil {
		return -1, err
	}

	// валидация токена
	if !token.Valid {
		return -1, errors.New("невалидный токен")
	}

	return claims.UserID, nil
}
