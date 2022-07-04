package router

import (
	"net/http"

	kecpsignal "github.com/fourdim/kecp/modules/kecp-signal"
	"github.com/fourdim/kecp/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func SetupKecpChiRouter() *chi.Mux {
	kecpRouter := chi.NewRouter()

	reg := kecpsignal.NewRegistry()

	kecpRouter.Route("/", func(r chi.Router) {
		r.Use(render.SetContentType(render.ContentTypeJSON))
		r.Post("/", services.NewRoomHandler(reg))
		r.Get("/", services.NewClientHandler(reg))
		r.Options("/", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
	})

	return kecpRouter
}
