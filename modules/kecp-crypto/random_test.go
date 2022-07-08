package kecpcrypto_test

import (
	"testing"

	. "github.com/fourdim/kecp/modules/kecp-crypto"
	"github.com/stretchr/testify/assert"
)

func TestGenerateCryptoKey(t *testing.T) {
	token := GenerateCryptoKey()
	assert.Len(t, token, 64)
}

func TestGenerateRoomID(t *testing.T) {
	token := GenerateRoomID()
	assert.Len(t, token, 16)
}
