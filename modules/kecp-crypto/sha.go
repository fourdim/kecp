package kecpcrypto

import "crypto/sha256"

func HashSha256(content []byte) []byte {
	hash := sha256.New()
	hash.Write(content)
	return hash.Sum(nil)
}
