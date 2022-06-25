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
	open     bool
	mx       sync.RWMutex
	reliable bool
	name     string
}

func NewConn(reliable bool, name string) *Conn {
	conn := &Conn{
		open:     true,
		reliable: reliable,
		name:     name,
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
	t := time.NewTimer(MathRandShortTimeGen())
	defer t.Stop()
	probability := MathRandGen()
	select {
	case <-t.C:
		if !conn.reliable && probability < 2 {
			return TextMessage, nil, &websocket.CloseError{
				Code: websocket.CloseAbnormalClosure,
				Text: "CloseAbnormalClosure",
			}
		} else {
			d, _ := json.Marshal(kecpmsg.Message{
				Type:    kecpmsg.Chat,
				Name:    conn.name,
				Payload: "hello",
			})
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
