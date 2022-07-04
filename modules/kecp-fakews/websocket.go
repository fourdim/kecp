package kecpfakews

import (
	"encoding/json"
	"errors"
	"io"
	"sync"
	"time"

	kecpmsg "github.com/fourdim/kecp/modules/kecp-msg"
	"github.com/gorilla/websocket"
)

var FakeError = errors.New("this is a fake error")

type Conn struct {
	open      bool
	auth      bool
	mx        sync.RWMutex
	reliable  bool
	roomID    string
	name      string
	clientKey string
}

func NewConn(reliable bool, roomID string, name string, clientKey string) *Conn {
	conn := &Conn{
		open:      true,
		auth:      false,
		reliable:  reliable,
		roomID:    roomID,
		name:      name,
		clientKey: clientKey,
	}
	return conn
}

func (conn *Conn) Close() error {
	conn.mx.Lock()
	conn.open = false
	conn.mx.Unlock()
	probability := MathRandGen()
	if !conn.reliable && probability < 2 {
		return FakeError
	} else {
		return nil
	}
}

func (conn *Conn) NextWriter(messageType int) (io.WriteCloser, error) {
	conn.mx.RLock()
	if !conn.open {
		conn.mx.RUnlock()
		return nil, &websocket.CloseError{
			Code: websocket.CloseNormalClosure,
			Text: "CloseNormalClosure",
		}
	}
	conn.mx.RUnlock()
	probability := MathRandGen()
	if !conn.reliable && probability < 2 {
		return nil, FakeError
	} else {
		return &FakeWriter{reliable: conn.reliable}, nil
	}
}

func (conn *Conn) SetPongHandler(h func(appData string) error) {
	t := time.NewTimer(MathRandShortTimeGen())
	defer t.Stop()
	select {
	case <-t.C:
		h("ping")
	}
}

func (conn *Conn) ReadMessage() (messageType int, p []byte, err error) {
	conn.mx.RLock()
	if !conn.open {
		conn.mx.RUnlock()
		return TextMessage, nil, &websocket.CloseError{
			Code: websocket.CloseNormalClosure,
			Text: "CloseNormalClosure",
		}
	}
	conn.mx.RUnlock()
	probability := MathRandGen()
	if !conn.auth {
		if !conn.reliable && probability < 2 {
			return TextMessage, []byte{}, nil
		}
		b, _ := json.Marshal(kecpmsg.AuthMessage{
			RoomID:    conn.roomID,
			Name:      conn.name,
			ClientKey: conn.clientKey,
		})
		conn.auth = true
		return TextMessage, b, nil
	}
	t := time.NewTimer(MathRandShortTimeGen())
	defer t.Stop()
	select {
	case <-t.C:
		if !conn.reliable && probability < 2 {
			return TextMessage, nil, &websocket.CloseError{
				Code: websocket.CloseAbnormalClosure,
				Text: "CloseAbnormalClosure",
			}
		} else {
			var msg interface{}

			if probability < 8 {
				msg = kecpmsg.Message{
					Type:    kecpmsg.Chat,
					Name:    conn.name,
					Payload: "hello",
				}
			} else if probability < 14 {
				msg = kecpmsg.Message{
					Type:    kecpmsg.Chat,
					Name:    conn.name,
					Target:  conn.name,
					Payload: "hello",
				}
			} else {
				msg = "error"
			}

			d, _ := json.Marshal(msg)
			return TextMessage, d, nil
		}
	}
}

func (conn *Conn) SetReadDeadline(t time.Time) error {
	probability := MathRandGen()
	if !conn.reliable && probability < 2 {
		return FakeError
	} else {
		return nil
	}
}

func (conn *Conn) SetReadLimit(limit int64) {}

func (conn *Conn) SetWriteDeadline(t time.Time) error {
	probability := MathRandGen()
	if !conn.reliable && probability < 2 {
		return FakeError
	} else {
		return nil
	}
}
func (conn *Conn) WriteMessage(messageType int, data []byte) error {
	probability := MathRandGen()
	if !conn.reliable && probability < 2 {
		return FakeError
	} else {
		return nil
	}
}
