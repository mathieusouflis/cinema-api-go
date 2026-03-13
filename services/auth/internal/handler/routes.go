package handler

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Munt(router *chi.Mux, dependencies *Dependencies) {

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "ok")
	})

	router.Route("/auth", func(r chi.Router) {
		r.Post("/login", NewLoginHandler(&dependencies.LoginUseCase).ServeHTTP)
		r.Post("/register", NewRegisterHandler(&dependencies.RegisterUseCase).ServeHTTP)
		r.Post("/refresh", NewRefreshHandler(&dependencies.RefreshUseCase).ServeHTTP)
		r.Post("/logout", NewLogoutHandler(&dependencies.LogoutUseCase).ServeHTTP)
		r.Post("/oauth/{provider}/callback", NewOauthHandler(&dependencies.OauthCallbackUseCase).ServeHTTP)
	})
}
