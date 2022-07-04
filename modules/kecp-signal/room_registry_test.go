package kecpsignal_test

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"math/rand"
	"net"
	"net/http"
	"testing"
	"time"

	kecpmsg "github.com/fourdim/kecp/modules/kecp-msg"
	. "github.com/fourdim/kecp/modules/kecp-signal"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

var upgrader = websocket.Upgrader{} // use default options

func echo(t *testing.T, reg *Registry, end chan bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("upgrade:", err)
			return
		}
		defer conn.Close()
		assert.NoError(t, reg.NewClient(conn))
		end <- true
	}
}

func TestGorillaWebsocketCompatibility(t *testing.T) {
	reg := NewRegistry()
	b := make([]byte, 48)
	rand.Read(b)
	userKey := base64.RawURLEncoding.EncodeToString(b)
	roomID := reg.NewRoom(userKey)

	end := make(chan bool)

	http.HandleFunc("/", echo(t, reg, end))
	l, err := net.Listen("tcp", "127.0.0.1:19216")
	assert.NoError(t, err, "error on listening to the port 19216")

	go http.Serve(l, http.DefaultServeMux)

	timer := time.NewTimer(1 * time.Millisecond)
	defer timer.Stop()
	select {
	case <-timer.C:
	}

	conn, _, err := websocket.DefaultDialer.Dial("ws://127.0.0.1:19216/", nil)
	auth, _ := json.Marshal(kecpmsg.AuthMessage{
		RoomID:    roomID,
		Name:      userKey[:16],
		ClientKey: userKey,
	})
	conn.WriteMessage(websocket.TextMessage, auth)
	assert.NoError(t, err, "error on dialing")
	<-end
}
