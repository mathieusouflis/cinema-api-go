package auth

import (
	"net/http"

	"example.com/filmserver/pkg/module"
)

func New(server *http.ServeMux) module.Module {
	authHandler := NewHandler()

	mod := module.New("/auth", server)

	mod.RegisterRoute(module.POST, "/login", http.HandlerFunc(authHandler.Login))
	mod.RegisterRoute(module.GET, "/register", http.HandlerFunc(authHandler.Register))

	return mod
}
