package service

import (
	"crypto/sha256"
	"encoding/hex"
)

const (
	hashKey = "asd52v04fgt2"
)

func genPasswordHash(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	hash.Write([]byte(hashKey))
	return hex.EncodeToString(hash.Sum(nil))
}
