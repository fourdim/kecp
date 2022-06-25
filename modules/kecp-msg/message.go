package kecpmsg

import (
	"encoding/json"
	"errors"
)

type MsgType string

type Message struct {
	// The message type.
	Type MsgType `json:"type"`

	// The sender's username.
	Name string `json:"name,omitempty"`

	// The username of the person to receive the message.
	Target string `json:"target,omitempty"`

	// The message's payload.
	Payload string `json:"payload"`

	// Broadcast except the client with clientKey.
	ExceptClientKey string `json:"-"`
}

const (
	VideoOffer      MsgType = "video-offer"
	VideoAnswer     MsgType = "video-answer"
	NewIceCandidate MsgType = "new-ice-candidate"
	Chat            MsgType = "chat"
	List            MsgType = "list"
	Join            MsgType = "join"
	Leave           MsgType = "leave"
)

var (
	ErrCanNotParseMessage = errors.New("can not prase the message")
)

func Parse(msg []byte, name string) (*Message, error) {
	var KecpMsg Message
	err := json.Unmarshal(msg, &KecpMsg)
	if err != nil {
		return nil, err
	}
	if KecpMsg.Name != name {
		return nil, ErrCanNotParseMessage
	}
	return &KecpMsg, nil
}

func (kecpMsg *Message) NeedBroadcast() bool {
	switch kecpMsg.Type {
	case VideoOffer:
		fallthrough
	case VideoAnswer:
		fallthrough
	case NewIceCandidate:
		fallthrough
	case List:
		return false
	case Chat:
		if kecpMsg.Target != "" {
			return false
		}
		fallthrough
	default:
		return true
	}
}

func (kecpMsg *Message) Build() []byte {
	b, _ := json.Marshal(kecpMsg)
	return b
}
