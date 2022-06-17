package kecpcrypto_test

import (
	"testing"

	. "github.com/fourdim/kecp/modules/kecp-crypto"
	"github.com/stretchr/testify/assert"
)

func TestGenerateToken(t *testing.T) {
	token := GenerateToken()
	assert.Len(t, token, 64)
}
