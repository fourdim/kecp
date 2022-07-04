package services

import (
	"net/http"

	kecpsignal "github.com/fourdim/kecp/modules/kecp-signal"
	"github.com/go-chi/render"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func NewClientHandler(reg *kecpsignal.Registry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			render.Render(w, r, ErrInvalidRequest(err))
			return
		}
		reg.NewClient(conn)
	}
}
