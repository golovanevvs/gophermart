package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"

	"github.com/golovanevvs/gophermart/internal/model"
	"github.com/golovanevvs/gophermart/internal/storage"
)

const (
	hashKey = "asd52v04fgt2"
)

type authServiceStr struct {
	st storage.AllStorageInt
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

func genPasswordHash(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	hash.Write([]byte(hashKey))
	return hex.EncodeToString(hash.Sum(nil))
}
