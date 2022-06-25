package kecpws_test

import (
	"encoding/base64"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"testing"
	"time"

	kecpfakews "github.com/fourdim/kecp/modules/kecp-fakews"
	. "github.com/fourdim/kecp/modules/kecp-signal"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

var upgrader = websocket.Upgrader{} // use default options

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func TestWithUsersJoined(t *testing.T) {
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

	for _, roomID := range roomIDs {
		for i := 0; i < 10; i++ {
			rand.Read(b)
			userKey := base64.RawURLEncoding.EncodeToString(b)
			userName := userKey[:12]
			reg.NewClient(userName, userKey, roomID, kecpfakews.NewConn(false, userName))
		}
	}

	timer1 := time.NewTimer(kecpfakews.MathRandLongTimeGen())
	defer timer1.Stop()
	select {
	case <-timer1.C:
	}

	for _, roomID := range roomIDs {
		if kecpfakews.MathRandGen() < 2 {
			reg.DeleteRoom(roomID)
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

	for _, roomID := range roomIDs {
		rand.Read(b)
		userKey := base64.RawURLEncoding.EncodeToString(b)
		userName := userKey[:12]
		for i := 0; i < 10; i++ {
			reg.NewClient(userName, userKey, roomID, kecpfakews.NewConn(false, userName))
		}
	}

	timer1 := time.NewTimer(kecpfakews.MathRandLongTimeGen())
	defer timer1.Stop()
	select {
	case <-timer1.C:
	}

	for _, roomID := range roomIDs {
		if kecpfakews.MathRandGen() < 2 {
			reg.DeleteRoom(roomID)
		}
	}

	timer2 := time.NewTimer(3 * time.Second)
	defer timer2.Stop()
	select {
	case <-timer2.C:
	}
}

func TestInvalidName(t *testing.T) {
	reg := NewRegistry()

	b := make([]byte, 48)

	rand.Read(b)
	userKey := base64.RawURLEncoding.EncodeToString(b)
	roomID, _ := reg.NewRoom(userKey)

	for i := 0; i < 3; i++ {
		rand.Read(b)
		userKey := base64.RawURLEncoding.EncodeToString(b)
		err := reg.NewClient(userKey, userKey, roomID, kecpfakews.NewConn(true, userKey))
		assert.EqualError(t, err, ErrNotAValidName.Error())
	}
}

func TestSameNames(t *testing.T) {
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

	for _, roomID := range roomIDs {
		for i := 0; i < 10; i++ {
			rand.Read(b)
			userKey := base64.RawURLEncoding.EncodeToString(b)
			err := reg.NewClient("〔=ヘ=#〕", userKey, roomID, kecpfakews.NewConn(true, "〔=ヘ=#〕"))
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

	for _, roomID := range roomIDs {
		if kecpfakews.MathRandGen() < 2 {
			reg.DeleteRoom(roomID)
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

	var roomIDs []string
	b := make([]byte, 48)
	for i := 0; i < 3; i++ {
		rand.Read(b)
		userKey := base64.RawURLEncoding.EncodeToString(b)
		roomID, _ := reg.NewRoom(userKey)
		if roomID != "" {
			roomIDs = append(roomIDs, roomID)
		}
	}

	for _, roomID := range roomIDs {
		for i := 0; i < 5; i++ {
			rand.Read(b)
			userKey := base64.RawURLEncoding.EncodeToString(b)
			userName := userKey[:12]
			reg.NewClient(userName, userKey, roomID, kecpfakews.NewConn(false, userName))
			if i == 2 {
				reg.DeleteRoom(roomID)
			}
		}
	}
}

func TestGoroutineLeak(t *testing.T) {
	reg := NewRegistry()

	baseGoroutineNum := runtime.NumGoroutine()

	var roomIDs []string
	b := make([]byte, 48)
	for i := 0; i < 5; i++ {
		rand.Read(b)
		userKey := base64.RawURLEncoding.EncodeToString(b)
		roomID, _ := reg.NewRoom(userKey)
		if roomID != "" {
			roomIDs = append(roomIDs, roomID)
		}
	}

	for _, roomID := range roomIDs {
		for i := 0; i < 3; i++ {
			rand.Read(b)
			userKey := base64.RawURLEncoding.EncodeToString(b)
			userName := userKey[:12]
			reg.NewClient(userName, userKey, roomID, kecpfakews.NewConn(true, userName))
		}
	}

	timer1 := time.NewTimer(kecpfakews.MathRandLongTimeGen())
	defer timer1.Stop()
	select {
	case <-timer1.C:
	}

	assert.GreaterOrEqual(t, runtime.NumGoroutine(), baseGoroutineNum+35)

	for _, roomID := range roomIDs {
		reg.DeleteRoom(roomID)
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
		roomID, _ := reg.NewRoom(userKey)
		if roomID != "" {
			roomIDs = append(roomIDs, roomID)
		}
	}

	for _, roomID := range roomIDs {
		for i := 0; i < 3; i++ {
			rand.Read(b)
			userKey := base64.RawURLEncoding.EncodeToString(b)
			userName := userKey[:12]
			reg.NewClient(userName, userKey, roomID, kecpfakews.NewConn(true, userName))
		}
	}

	timer1 := time.NewTimer(61 * time.Second)
	defer timer1.Stop()
	select {
	case <-timer1.C:
	}
}
