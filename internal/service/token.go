package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/golovanevvs/gophermart/internal/customerrors"
)

const (
	TokenExp  = time.Hour * 3 // время жизни токена
	SecretKey = "sskey"       // секретный ключ токена
)

// структура утверждений
type claims struct {
	jwt.RegisteredClaims
	UserID int
}

// BuildJWTString создаёт токен и возвращает его в виде строки
func (as *authServiceStr) BuildJWTString(ctx context.Context, login, password string) (string, error) {
	// получение пользователя из БД
	user, err := as.st.LoadUserByLoginPasswordHash(ctx, login, genPasswordHash(password))
	if err != nil {
		// если неверная пара логин/пароль
		if strings.Contains(err.Error(), "no rows in result set") {
			return "", fmt.Errorf("%v: %v", customerrors.DBInvalidLoginPassword401, err.Error())
		}
		// если другая ошибка
		return "", fmt.Errorf("%v: %v", customerrors.InternalServerError500, err.Error())
	}

	// создание токена
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExp)),
		},
		UserID: user.ID,
	})

	// создание строки токена
	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return "", fmt.Errorf("%v: %v", customerrors.InternalServerError500, err.Error())
	}
	return tokenString, nil
}

// GetUserIDFromJWT возвращает userID из JWT
func (as *authServiceStr) GetUserIDFromJWT(tokenString string) (int, error) {
	claims := &claims{}

	// преобразование строки в токен
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("неверный метод подписи")
		}
		return []byte(SecretKey), nil
	})
	if err != nil {
		return -1, fmt.Errorf("%v: %v", customerrors.JWTParseError401, err.Error())
	}

	// валидация токена
	if !token.Valid {
		return -1, fmt.Errorf("%v", customerrors.JWTInvalidToken401)
	}

	return claims.UserID, nil
}
