package service

import (
	"fmt"

	"github.com/golovanevvs/gophermart/internal/storage"
)

type AuthServiceStr struct {
	st storage.AuthInt
}

func NewAuthService(st storage.AuthInt) *AuthServiceStr {
	return &AuthServiceStr{
		st: st,
	}
}

func (as *AuthServiceStr) CreateUser() string {
	return fmt.Sprintf("CreateUser")
}
