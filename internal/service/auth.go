package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/golovanevvs/gophermart/internal/model"
	"github.com/golovanevvs/gophermart/internal/storage"
)

const (
	hashKey    = "asd52v04fgt2"
	TOKEN_EXP  = time.Hour * 3
	SECRET_KEY = "sskey"
)

type authServiceStr struct {
	st storage.AllStorageInt
}

type claims struct {
	jwt.RegisteredClaims
	UserID int
}

func NewAuthService(st storage.AllStorageInt) *authServiceStr {
	return &authServiceStr{
		st: st,
	}
}

func (as *authServiceStr) CreateUser(ctx context.Context, user model.User) error {
	user.Password = genPasswordHash(user.Password)
	return as.st.CreateUser(ctx, user)
}

func (as *authServiceStr) GenToken(ctx context.Context, login, password string) (string, error) {
	user, err := as.st.GetUser(ctx, login, genPasswordHash(password))
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

func genPasswordHash(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	hash.Write([]byte(hashKey))
	return hex.EncodeToString(hash.Sum(nil))
}
