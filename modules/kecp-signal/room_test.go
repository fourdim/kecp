package kecpws_test

import (
	"encoding/base64"
	"math/rand"
	"testing"

	kecpfakews "github.com/fourdim/kecp/modules/kecp-fakews"
	. "github.com/fourdim/kecp/modules/kecp-signal"
)

func TestNoUserJoined(t *testing.T) {
	reg := NewRegistry()

	var roomIDs []string
	b := make([]byte, 48)
	for i := 0; i < 10; i++ {
		rand.Read(b)
		userKey := base64.RawURLEncoding.EncodeToString(b)
		roomID, _ := reg.NewRoom(userKey)
		if roomID != "" {
			roomIDs = append(roomIDs, roomID)
		}
	}

	for range roomIDs {
		rand.Read(b)
		userKey := base64.RawURLEncoding.EncodeToString(b)
		rand.Read(b)
		roomID := base64.RawURLEncoding.EncodeToString(b)
		reg.NewClient(userKey, roomID, kecpfakews.NewConn(false))
	}
}
