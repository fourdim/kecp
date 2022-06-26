package services

import (
	"errors"
	"net/http"

	kecpsignal "github.com/fourdim/kecp/modules/kecp-signal"
	kecpvalidate "github.com/fourdim/kecp/modules/kecp-validate"
	"github.com/go-chi/render"
)

type CreateRoomRequest struct {
	ClientKey string `json:"client_key"`
}

func (req *CreateRoomRequest) Bind(r *http.Request) error {
	if !kecpvalidate.IsAValidCryptoKey(req.ClientKey) {
		return errors.New("malformed client key.")
	}
	return nil
}

type CreateRoomResponse struct {
	RoomID string `json:"room_id"`
}

func (resp *CreateRoomResponse) Render(w http.ResponseWriter, r *http.Request) error {
	if len(resp.RoomID) != 64 {
		return errors.New("unable to create the room.")
	}
	return nil
}

func NewRoomHandler(reg *kecpsignal.Registry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &CreateRoomRequest{}
		if err := render.Bind(r, req); err != nil {
			render.Render(w, r, ErrInvalidRequest(err))
			return
		}
		roomID := reg.NewRoom(req.ClientKey)
		resp := &CreateRoomResponse{RoomID: roomID}
		if err := render.Render(w, r, resp); err != nil {
			render.Render(w, r, ErrInternalError(err))
			return
		}
	}
}
