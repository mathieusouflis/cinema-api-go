package auth

import (
	"net/http"

	router "example.com/filmserver/pkg/module"
)

func New(server *http.ServeMux) router.Router {
	authHandler := NewHandler()

	mod := router.New("/auth", server)

	mod.RegisterRoute(router.POST, "/login", http.HandlerFunc(authHandler.Login))
	mod.RegisterRoute(router.GET, "/register", http.HandlerFunc(authHandler.Register))

	return mod
}
