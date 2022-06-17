package kecpws

import (
	"bytes"
	"io"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512

	// Time allowed to get an ack from a room.
	clientJoinedCheckWait = 2 * time.Second
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type Client struct {
	ClientKey string

	// The room it belongs.
	room *Room

	// The websocket connection.
	conn WebscoketConn

	// Buffered channel of outbound messages.
	send chan []byte

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

func (reg *Registry) NewClient(key string, roomID string, conn WebscoketConn) error {
	room := reg.getRoom(roomID)
	if room == nil {
		conn.Close()
		return ErrCanNotJoinTheRoom
	}
	client := &Client{
		ClientKey:       key,
		room:            room,
		conn:            conn,
		send:            make(chan []byte, 256),
		joined:          make(chan bool),
		selfDestruction: make(chan bool),
	}
	client.room.register.Write(client)
	checker := time.NewTimer(clientJoinedCheckWait)
	defer checker.Stop()
	select {
	case <-client.joined:
	case <-checker.C:
		client.conn.Close()
		return ErrCanNotJoinTheRoom
	}
	go client.readPump()
	go client.writePump()
	return nil
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
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logger.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		c.room.broadcast.Write(message)
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
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The room closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		case <-c.selfDestruction:
			return
		}
	}
}
