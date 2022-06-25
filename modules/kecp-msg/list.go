package kecpmsg

import (
	"encoding/json"
)

func NewListMessage(list []string) *Message {
	p, _ := json.Marshal(list)
	return &Message{
		Type:    List,
		Payload: string(p),
	}
}

func NewJoinMessage(name string, clientKey string) *Message {
	return &Message{
		Type:            Join,
		Payload:         name,
		ExceptClientKey: clientKey,
	}
}

func NewLeaveMessage(name string) *Message {
	return &Message{
		Type:    Join,
		Payload: name,
	}
}
