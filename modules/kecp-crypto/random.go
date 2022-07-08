package kecpcrypto

import (
	"crypto/rand"
	"encoding/base64"
)

// GenerateCryptoKey generate a string of 64 characters.
func GenerateCryptoKey() (token string) {
	b := make([]byte, 48)
	rand.Read(b)
	token = base64.RawURLEncoding.EncodeToString(b)
	return
}

func GenerateRoomID() (roomID string) {
	b := make([]byte, 12)
	rand.Read(b)
	roomID = base64.RawURLEncoding.EncodeToString(b)
	return
}
