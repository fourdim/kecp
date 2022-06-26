package services

import (
	"errors"
	"net/http"

	kecpsignal "github.com/fourdim/kecp/modules/kecp-signal"
	kecpvalidate "github.com/fourdim/kecp/modules/kecp-validate"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func NewClientHandler(reg *kecpsignal.Registry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		roomID := chi.URLParam(r, "roomID")
		query := r.URL.Query()
		name := query.Get("name")
		clientKey := query.Get("client_key")
		if !kecpvalidate.IsAValidUserName(name) {
			render.Render(w, r, ErrInvalidRequest(errors.New("invalid user name.")))
			return
		}
		if !kecpvalidate.IsAValidUserName(clientKey) {
			render.Render(w, r, ErrInvalidRequest(errors.New("invalid client key length or too weak.")))
			return
		}
		room := reg.GetRoom(roomID)
		if room == nil {
			render.Render(w, r, ErrInvalidRequest(errors.New("invalid room id.")))
			return
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			render.Render(w, r, ErrInvalidRequest(err))
			return
		}
		room.NewClient(name, clientKey, conn)
	}
}
