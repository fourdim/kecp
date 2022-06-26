package kecpsignal_test

import (
	"net"
	"net/http"
	"testing"
	"time"

	. "github.com/fourdim/kecp/modules/kecp-signal"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestGorillaWebsocketCompatibility(t *testing.T) {
	reg := NewRegistry()
	roomID := reg.NewRoom("aaa")

	http.HandleFunc("/", echo)
	l, err := net.Listen("tcp", "127.0.0.1:19216")
	assert.NoError(t, err, "error on listening to the port 19216")

	go http.Serve(l, http.DefaultServeMux)

	timer := time.NewTimer(1 * time.Second)
	defer timer.Stop()
	select {
	case <-timer.C:
	}

	con, _, err := websocket.DefaultDialer.Dial("ws://127.0.0.1:19216/", nil)
	assert.NoError(t, err, "error on dialing")
	assert.NoError(t, reg.GetRoom(roomID).NewClient("〔=ヘ=#〕", "aaa", con))
}
