package kecpmsg_test

import (
	"testing"

	. "github.com/fourdim/kecp/modules/kecp-msg"
	"github.com/stretchr/testify/assert"
)

func TestParseVideoOfferMessage(t *testing.T) {
	msg, err := Parse([]byte(`{"type":"video-offer","name":"Alice","target":"Bob","payload":"Hello"}`), "Alice")
	assert.NoError(t, err)
	assert.False(t, msg.NeedBroadcast())
	assert.Equal(t, `{"type":"video-offer","name":"Alice","target":"Bob","payload":"Hello"}`, string(msg.Build()))
}

func TestParseVideoAnswerMessage(t *testing.T) {
	msg, err := Parse([]byte(`{"type":"video-offer","name":"Bob","target":"Alice","payload":"Hello"}`), "Bob")
	assert.NoError(t, err)
	assert.False(t, msg.NeedBroadcast())
	assert.Equal(t, `{"type":"video-offer","name":"Bob","target":"Alice","payload":"Hello"}`, string(msg.Build()))
}

func TestParseNewIceCandidateMessage(t *testing.T) {
	msg, err := Parse([]byte(`{"type":"new-ice-candidate","name":"Alice","target":"Bob","payload":"Hello"}`), "Alice")
	assert.NoError(t, err)
	assert.False(t, msg.NeedBroadcast())
	assert.Equal(t, `{"type":"new-ice-candidate","name":"Alice","target":"Bob","payload":"Hello"}`, string(msg.Build()))
}

func TestParseChatToAllMessage(t *testing.T) {
	msg, err := Parse([]byte(`{"type":"chat","name":"Alice","payload":"Hello"}`), "Alice")
	assert.NoError(t, err)
	assert.True(t, msg.NeedBroadcast())
	assert.Equal(t, `{"type":"chat","name":"Alice","payload":"Hello"}`, string(msg.Build()))
}

func TestParseChatToOneMessage(t *testing.T) {
	msg, err := Parse([]byte(`{"type":"chat","name":"Alice","target":"Bob","payload":"Hello"}`), "Alice")
	assert.NoError(t, err)
	assert.False(t, msg.NeedBroadcast())
	assert.Equal(t, `{"type":"chat","name":"Alice","target":"Bob","payload":"Hello"}`, string(msg.Build()))
}

func TestParseListMessage(t *testing.T) {
	_, err := Parse([]byte(`{"type":"list","name":"Mallory","payload":["Alice","Bob"]}`), "Mallory")
	assert.EqualError(t, err, ErrCanNotParseMessage.Error())
}

func TestParseJoinMessage(t *testing.T) {
	_, err := Parse([]byte(`{"type":"join","name":"Mallory","payload":"Alice"}`), "Mallory")
	assert.EqualError(t, err, ErrCanNotParseMessage.Error())
}

func TestParseLeaveMessage(t *testing.T) {
	_, err := Parse([]byte(`{"type":"leave","name":"Mallory","payload":"Alice"}`), "Mallory")
	assert.EqualError(t, err, ErrCanNotParseMessage.Error())
}

func TestParseMalformedMessage(t *testing.T) {
	_, err := Parse([]byte(`{"type":"chat","name":"Mallory","target":"Bob","payload:"Hello"}`), "Mallory")
	assert.EqualError(t, err, "invalid character 'H' after object key")
}

func TestParseFraudMessage(t *testing.T) {
	_, err := Parse([]byte(`{"type":"chat","name":"Alice","target":"Bob","payload":"Hello"}`), "Mallory")
	assert.EqualError(t, err, ErrCanNotParseMessage.Error())
}
