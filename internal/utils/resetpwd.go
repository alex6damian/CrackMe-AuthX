package utils

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateSecureResetToken() string {
	bytes := make([]byte, 32)
	_, _ = rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
