package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

func GenerateSecureString(length int) (string, string) {
	bytes := make([]byte, length)
	rand.Read(bytes)

	plainToken := hex.EncodeToString(bytes)[:length]
	hash := sha256.Sum256([]byte(plainToken))
	hashedToken := hex.EncodeToString(hash[:])

	return plainToken, hashedToken
}
