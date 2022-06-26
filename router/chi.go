package router

import (
	kecpsignal "github.com/fourdim/kecp/modules/kecp-signal"
	"github.com/fourdim/kecp/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func SetupKecpChiRouter() *chi.Mux {
	kecpRouter := chi.NewRouter()

	reg := kecpsignal.NewRegistry()

	kecpRouter.Route("/rooms", func(r chi.Router) {
		r.Use(render.SetContentType(render.ContentTypeJSON))
		r.Post("/", services.NewRoomHandler(reg))
		r.Get("/{roomID}", services.NewClientHandler(reg))
	})

	return kecpRouter
}
