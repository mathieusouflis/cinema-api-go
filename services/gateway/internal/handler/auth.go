package handler

import (
	"gateway/internal/config"
	reverseproxy "gateway/pkg/reverse-proxy"

	"github.com/go-chi/chi/v5"
)

func NewAuthHandler(router *chi.Mux) {
	cfg := config.Load()

	proxy := reverseproxy.New(cfg.AuthServiceURL)

	router.Route("/auth", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Mount("/login", proxy)
			r.Mount("/register", proxy)
			r.Mount("/logout", proxy)
			r.Mount("/refresh", proxy)
			r.Mount("/oauth/{provider}/callback", proxy)
		})
	})

}
