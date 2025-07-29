package auth

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateRandomToken(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}