package kecpws_test

import (
	"encoding/base64"
	"math/rand"
	"testing"
	"time"

	kecpfakews "github.com/fourdim/kecp/modules/kecp-fakews"
	. "github.com/fourdim/kecp/modules/kecp-signal"
	"github.com/stretchr/testify/assert"
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
		err := reg.NewClient("〔=ヘ=#〕", userKey, userKey, kecpfakews.NewConn(false, "〔=ヘ=#〕"))
		assert.EqualError(t, err, ErrCanNotJoinTheRoom.Error())
	}

	timer1 := time.NewTimer(31 * time.Second)
	defer timer1.Stop()
	select {
	case <-timer1.C:
	}
}
