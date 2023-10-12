package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

func HashPassword(password string) string {
	hashPass := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hashPass[:])
}
