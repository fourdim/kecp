package kecpsignal

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"time"

	kecpmsg "github.com/fourdim/kecp/modules/kecp-msg"
	kecpvalidate "github.com/fourdim/kecp/modules/kecp-validate"
	ws "github.com/gorilla/websocket"
)

const (
	// Time allowed to read the auth message.
	authWait = 5 * time.Second

	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 10240

	// Time allowed to get an ack from a room.
	clientJoinedCheckWait = 2 * time.Second
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type Client struct {

	// The clientKey is generated by the client.
	// Should be readonly.
	clientKey string

	// The name is provided by the client.
	// Should be unique in one room, and will be known to other clients.
	// Should be readonly.
	name string

	// The room it belongs.
	room *Room

	// The websocket connection.
	conn WebscoketConn

	// Buffered channel of outbound messages.
	send chan *kecpmsg.Message

	// The status returned after register.
	joined chan bool

	// Channel for self destruction.
	selfDestruction chan bool
}

type WebscoketConn interface {
	Close() error
	NextWriter(messageType int) (io.WriteCloser, error)
	SetPongHandler(h func(appData string) error)
	ReadMessage() (messageType int, p []byte, err error)
	SetReadDeadline(t time.Time) error
	SetReadLimit(limit int64)
	SetWriteDeadline(t time.Time) error
	WriteMessage(messageType int, data []byte) error
}

func (reg *Registry) NewClient(conn WebscoketConn) (retErr error) {
	defer func() {
		if retErr != nil && !errors.Is(retErr, ErrConnectionLost) {
			sendErrorMsg(conn, retErr)
		}
		if retErr != nil {
			conn.Close()
		}
	}()
	conn.SetReadLimit(maxMessageSize)
	conn.SetReadDeadline(time.Now().Add(authWait))
	_, msg, err := conn.ReadMessage()
	if err != nil {
		return ErrConnectionLost
	}
	var auth kecpmsg.AuthMessage
	conn.SetWriteDeadline(time.Now().Add(writeWait))
	if err := json.Unmarshal(msg, &auth); err != nil {
		return ErrCanNotJoinTheRoom
	}
	if !kecpvalidate.IsAValidCryptoKey(auth.ClientKey) {
		return ErrNotAValidKey
	}
	if !kecpvalidate.IsAValidUserName(auth.Name) {
		return ErrNotAValidName
	}
	room := reg.GetRoom(auth.RoomID)
	if room == nil {
		return ErrCanNotJoinTheRoom
	}
	client := &Client{
		clientKey:       auth.ClientKey,
		name:            auth.Name,
		room:            room,
		conn:            conn,
		send:            make(chan *kecpmsg.Message, 256),
		joined:          make(chan bool),
		selfDestruction: make(chan bool),
	}
	room.register.Write(client)
	checker := time.NewTimer(clientJoinedCheckWait)
	defer checker.Stop()
	select {
	case joined := <-client.joined:
		if !joined {
			return ErrNameIsAlreadyInUse
		}
	case <-checker.C:
		sendErrorMsg(conn, ErrCanNotJoinTheRoom)
		return ErrCanNotJoinTheRoom
	}
	client.sendListMsg()
	go client.readPump()
	go client.writePump()
	return nil
}

func sendErrorMsg(conn WebscoketConn, err error) {
	conn.WriteMessage(ws.TextMessage, kecpmsg.NewErrorMsg(err).Build())
	conn.WriteMessage(ws.CloseMessage, []byte{})
}

func (c *Client) sendListMsg() {
	c.conn.WriteMessage(ws.TextMessage, (<-c.send).Build())
}

// readPump pumps messages from the websocket connection to the room.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
//
// server <- client
func (c *Client) readPump() {
	defer func() {
		c.room.unregister.Write(c)
		c.conn.Close()
	}()
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			if ws.IsUnexpectedCloseError(err, ws.CloseGoingAway, ws.CloseAbnormalClosure) {
				logger.Printf("error: %v", err)
			}
			break
		}
		msg = bytes.TrimSpace(bytes.Replace(msg, newline, space, -1))
		kecpMsg, err := kecpmsg.Parse(msg, c.name)
		if err != nil {
			continue
		}

		if kecpMsg.NeedBroadcast() {
			c.room.broadcast.Write(kecpMsg)
		} else {
			c.room.forward.Write(kecpMsg)
		}
	}
}

// writePump pumps messages from the room to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
//
// server -> client
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case kecpMsg, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The room closed the channel.
				c.conn.WriteMessage(ws.CloseMessage, []byte{})
				return
			}
			err := c.conn.WriteMessage(ws.TextMessage, kecpMsg.Build())
			if err != nil {
				return
			}
			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				kecpMsg := <-c.send
				err := c.conn.WriteMessage(ws.TextMessage, kecpMsg.Build())
				if err != nil {
					return
				}
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(ws.PingMessage, nil); err != nil {
				return
			}
		case <-c.selfDestruction:
			return
		}
	}
}
