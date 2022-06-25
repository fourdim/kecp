package router

import (
	"github.com/go-chi/chi/v5"
)

func SetupChiRouter() *chi.Mux {
	kecpRouter := chi.NewRouter()

	kecpRouter.Route("/reg", func(r chi.Router) {

	})

	kecpRouter.Route("/ws", func(r chi.Router) {

	})

	return kecpRouter
}
