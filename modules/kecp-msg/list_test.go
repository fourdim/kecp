package kecpmsg_test

import (
	"testing"

	. "github.com/fourdim/kecp/modules/kecp-msg"
	"github.com/stretchr/testify/assert"
)

func TestNewListMessage(t *testing.T) {
	msg := NewListMessage([]string{"Alice", "Bob"})
	assert.Equal(t, "[\"Alice\",\"Bob\"]", msg.Payload)
}
