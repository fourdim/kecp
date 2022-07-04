package kecpsignal_test

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
		roomID := reg.NewRoom(userKey)
		if roomID != "" {
			roomIDs = append(roomIDs, roomID)
		}
	}

	for range roomIDs {
		rand.Read(b)
		userKey := base64.RawURLEncoding.EncodeToString(b)
		err := reg.NewClient(kecpfakews.NewConn(false, userKey, "〔=ヘ=#〕", userKey))
		assert.EqualError(t, err, ErrCanNotJoinTheRoom.Error())
	}

	timer1 := time.NewTimer(31 * time.Second)
	defer timer1.Stop()
	select {
	case <-timer1.C:
	}

	for _, roomID := range roomIDs {
		rand.Read(b)
		userKey := base64.RawURLEncoding.EncodeToString(b)
		userName := userKey[:12]
		assert.EqualError(t, reg.NewClient(kecpfakews.NewConn(true, roomID, userName, userKey)), ErrCanNotJoinTheRoom.Error())
	}
}

func TestUserReplacement(t *testing.T) {
	reg := NewRegistry()
	b := make([]byte, 48)
	rand.Read(b)
	userKey := base64.RawURLEncoding.EncodeToString(b)
	roomID := reg.NewRoom(userKey)
	err := reg.NewClient(kecpfakews.NewConn(true, roomID, "〔=ヘ=#〕", userKey))
	assert.NoError(t, err)

	timer1 := time.NewTimer(1 * time.Millisecond)
	defer timer1.Stop()
	select {
	case <-timer1.C:
	}

	err = reg.NewClient(kecpfakews.NewConn(true, roomID, "〔=ヘ=#〕", userKey))
	assert.NoError(t, err)
}
