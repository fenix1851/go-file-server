package utils

import (
	"crypto/sha256"
)

func HashPassword(password string) [32]byte {
	return sha256.Sum256([]byte(password))
}
