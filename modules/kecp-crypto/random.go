package kecpcrypto

import (
	"crypto/rand"
	"encoding/base64"
)

// GenerateToken generate a string of 64 characters.
func GenerateToken() (token string) {
	b := make([]byte, 48)
	rand.Read(b)
	token = base64.RawURLEncoding.EncodeToString(b)
	return
}
