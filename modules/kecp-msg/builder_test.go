package kecpmsg_test

import (
	"errors"
	"testing"

	. "github.com/fourdim/kecp/modules/kecp-msg"
	"github.com/stretchr/testify/assert"
)

func TestNewListMessage(t *testing.T) {
	msg := NewListMsg([]string{"Alice", "Bob"})
	assert.Equal(t, `{"type":"list","payload":["Alice","Bob"]}`, string(msg.Build()))
}

func TestNewJoinMessage(t *testing.T) {
	msg := NewJoinMsg("Alice", "aaa")
	assert.Equal(t, `{"type":"join","payload":"Alice"}`, string(msg.Build()))
}

func TestNewLeaveMessage(t *testing.T) {
	msg := NewLeaveMsg("Alice", "aaa")
	assert.Equal(t, `{"type":"leave","payload":"Alice"}`, string(msg.Build()))
}

func TestNewErrorMessage(t *testing.T) {
	msg := NewErrorMsg(errors.New("err"))
	assert.Equal(t, `{"type":"error","payload":"err"}`, string(msg.Build()))
}
