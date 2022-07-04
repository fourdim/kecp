package kecpsignal_test

import (
	"encoding/base64"
	"math/rand"
	"runtime"
	"testing"
	"time"

	kecpfakews "github.com/fourdim/kecp/modules/kecp-fakews"
	. "github.com/fourdim/kecp/modules/kecp-signal"
	"github.com/stretchr/testify/assert"
)

type FakeRoom struct {
	roomID        string
	managementKey string
}

func TestWithUsersJoined(t *testing.T) {
	reg := NewRegistry()

	var rooms []*FakeRoom
	b := make([]byte, 48)
	for i := 0; i < 10; i++ {
		rand.Read(b)
		userKey := base64.RawURLEncoding.EncodeToString(b)
		roomID := reg.NewRoom(userKey)
		if roomID != "" {
			rooms = append(rooms, &FakeRoom{roomID: roomID, managementKey: userKey})
		}
	}

	for _, room := range rooms {
		for i := 0; i < 10; i++ {
			rand.Read(b)
			userKey := base64.RawURLEncoding.EncodeToString(b)
			userName := userKey[:12]
			reg.NewClient(kecpfakews.NewConn(false, room.roomID, userName, userKey))
		}
	}

	timer1 := time.NewTimer(kecpfakews.MathRandLongTimeGen())
	defer timer1.Stop()
	select {
	case <-timer1.C:
	}

	for _, room := range rooms {
		if kecpfakews.MathRandGen() < 2 {
			reg.DeleteRoom(room.roomID, room.managementKey)
		}
	}

	timer2 := time.NewTimer(3 * time.Second)
	defer timer2.Stop()
	select {
	case <-timer2.C:
	}
}

func TestSameUsers(t *testing.T) {
	reg := NewRegistry()

	var rooms []*FakeRoom
	b := make([]byte, 48)
	for i := 0; i < 10; i++ {
		rand.Read(b)
		userKey := base64.RawURLEncoding.EncodeToString(b)
		roomID := reg.NewRoom(userKey)
		if roomID != "" {
			rooms = append(rooms, &FakeRoom{roomID: roomID, managementKey: userKey})
		}
	}

	for _, room := range rooms {
		rand.Read(b)
		userKey := base64.RawURLEncoding.EncodeToString(b)
		userName := userKey[:12]
		for i := 0; i < 10; i++ {
			reg.NewClient(kecpfakews.NewConn(false, room.roomID, userName, userKey))
		}
	}
}

func TestInvalidName(t *testing.T) {
	reg := NewRegistry()

	b := make([]byte, 48)

	rand.Read(b)
	userKey := base64.RawURLEncoding.EncodeToString(b)
	roomID := reg.NewRoom(userKey)

	for i := 0; i < 3; i++ {
		rand.Read(b)
		userKey := base64.RawURLEncoding.EncodeToString(b)
		err := reg.NewClient(kecpfakews.NewConn(true, roomID, userKey, userKey))
		assert.EqualError(t, err, ErrNotAValidName.Error())
	}
}

func TestInvalidKey(t *testing.T) {
	reg := NewRegistry()

	b := make([]byte, 48)

	rand.Read(b)
	userKey := base64.RawURLEncoding.EncodeToString(b)
	roomID := reg.NewRoom(userKey)

	for i := 0; i < 3; i++ {
		rand.Read(b)
		userKey := base64.RawURLEncoding.EncodeToString(b)
		err := reg.NewClient(kecpfakews.NewConn(true, roomID, userKey, userKey[:16]))
		assert.EqualError(t, err, ErrNotAValidKey.Error())
	}
}

func TestSameNames(t *testing.T) {
	reg := NewRegistry()

	var rooms []*FakeRoom
	b := make([]byte, 48)
	for i := 0; i < 10; i++ {
		rand.Read(b)
		userKey := base64.RawURLEncoding.EncodeToString(b)
		roomID := reg.NewRoom(userKey)
		if roomID != "" {
			rooms = append(rooms, &FakeRoom{roomID: roomID, managementKey: userKey})
		}
	}

	for _, room := range rooms {
		for i := 0; i < 10; i++ {
			rand.Read(b)
			userKey := base64.RawURLEncoding.EncodeToString(b)
			err := reg.NewClient(kecpfakews.NewConn(true, room.roomID, "〔=ヘ=#〕", userKey))
			if i != 0 {
				assert.EqualError(t, err, ErrNameIsAlreadyInUse.Error())
			}
		}
	}

	timer1 := time.NewTimer(kecpfakews.MathRandLongTimeGen())
	defer timer1.Stop()
	select {
	case <-timer1.C:
	}

	for _, room := range rooms {
		if kecpfakews.MathRandGen() < 2 {
			reg.DeleteRoom(room.roomID, room.managementKey)
		}
	}

	timer2 := time.NewTimer(3 * time.Second)
	defer timer2.Stop()
	select {
	case <-timer2.C:
	}
}

func TestWithUsersJoinedAndRoomDeletedSimultaneously(t *testing.T) {
	reg := NewRegistry()

	var rooms []*FakeRoom
	b := make([]byte, 48)
	for i := 0; i < 3; i++ {
		rand.Read(b)
		userKey := base64.RawURLEncoding.EncodeToString(b)
		roomID := reg.NewRoom(userKey)
		if roomID != "" {
			rooms = append(rooms, &FakeRoom{roomID: roomID, managementKey: userKey})
		}
	}

	for _, room := range rooms {
		for i := 0; i < 5; i++ {
			rand.Read(b)
			userKey := base64.RawURLEncoding.EncodeToString(b)
			userName := userKey[:12]
			reg.NewClient(kecpfakews.NewConn(false, room.roomID, userName, userKey))
			if i == 2 {
				reg.DeleteRoom(room.roomID, room.managementKey)
			}
		}
	}
}

func TestGoroutineLeak(t *testing.T) {
	reg := NewRegistry()

	baseGoroutineNum := runtime.NumGoroutine()

	var rooms []*FakeRoom
	b := make([]byte, 48)
	for i := 0; i < 5; i++ {
		rand.Read(b)
		userKey := base64.RawURLEncoding.EncodeToString(b)
		roomID := reg.NewRoom(userKey)
		if roomID != "" {
			rooms = append(rooms, &FakeRoom{roomID: roomID, managementKey: userKey})
		}
	}

	for _, room := range rooms {
		for i := 0; i < 3; i++ {
			rand.Read(b)
			userKey := base64.RawURLEncoding.EncodeToString(b)
			userName := userKey[:12]
			assert.NoError(t, reg.NewClient(kecpfakews.NewConn(true, room.roomID, userName, userKey)))
		}
	}

	timer1 := time.NewTimer(kecpfakews.MathRandLongTimeGen())
	defer timer1.Stop()
	select {
	case <-timer1.C:
	}

	// Add a breakpoint here to see whether goroutine leaks.
	for _, room := range rooms {
		reg.DeleteRoom(room.roomID, room.managementKey)
	}

	timer2 := time.NewTimer(3 * time.Second)
	defer timer2.Stop()
	select {
	case <-timer2.C:
	}

	// Add a breakpoint here to see whether goroutine leaks.
	assert.LessOrEqual(t, runtime.NumGoroutine(), baseGoroutineNum+6)
}

func TestClientPing(t *testing.T) {
	reg := NewRegistry()

	var roomIDs []string
	b := make([]byte, 48)
	for i := 0; i < 5; i++ {
		rand.Read(b)
		userKey := base64.RawURLEncoding.EncodeToString(b)
		roomID := reg.NewRoom(userKey)
		if roomID != "" {
			roomIDs = append(roomIDs, roomID)
		}
	}

	for _, roomID := range roomIDs {
		for i := 0; i < 3; i++ {
			rand.Read(b)
			userKey := base64.RawURLEncoding.EncodeToString(b)
			userName := userKey[:12]
			assert.NoError(t, reg.NewClient(kecpfakews.NewConn(true, roomID, userName, userKey)))
		}
	}

	timer1 := time.NewTimer(61 * time.Second)
	defer timer1.Stop()
	select {
	case <-timer1.C:
	}
}
