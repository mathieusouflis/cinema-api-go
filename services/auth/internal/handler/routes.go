package handler

import (
	"github.com/go-chi/chi/v5"
)

func Munt(router *chi.Mux, dependencies *Dependencies) {
	router.Route("/auth", func(r chi.Router) {
		r.Post("/login", NewLoginHandler(&dependencies.LoginUseCase).ServeHTTP)
	})
}
