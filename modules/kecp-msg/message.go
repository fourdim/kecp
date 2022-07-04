package kecpmsg

import (
	"encoding/json"
	"errors"
	"log"
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
	Payload interface{} `json:"payload"`

	// Broadcast except the client with clientKey.
	ExceptClientKey string `json:"-"`
}

type AuthMessage struct {
	RoomID    string `json:"room_id"`
	Name      string `json:"name"`
	ClientKey string `json:"client_key"`
}

const (
	VideoOffer      MsgType = "video-offer"
	VideoAnswer     MsgType = "video-answer"
	NewIceCandidate MsgType = "new-ice-candidate"
	Chat            MsgType = "chat"
	List            MsgType = "list"
	Join            MsgType = "join"
	Leave           MsgType = "leave"
	Error           MsgType = "error"
)

var (
	ErrCanNotParseMessage = errors.New("can not prase the message")
)

func Parse(msg []byte, name string) (*Message, error) {
	var kecpMsg Message
	err := json.Unmarshal(msg, &kecpMsg)
	if err != nil {
		return nil, err
	}
	if kecpMsg.Name != name {
		return nil, ErrCanNotParseMessage
	}
	switch kecpMsg.Type {
	case List:
		fallthrough
	case Join:
		fallthrough
	case Leave:
		return nil, ErrCanNotParseMessage
	}
	return &kecpMsg, nil
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
		fallthrough
	case Join:
		fallthrough
	case Leave:
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
	log.Println(string(b))
	return b
}
