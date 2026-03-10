package films

import (
	"net/http"

	router "example.com/filmserver/pkg/module"
)

func New(server *http.ServeMux) router.Router {
	filmsHandler := NewHandler()

	mod := router.New("/films", server)

	mod.RegisterRoute(router.POST, "/list", http.HandlerFunc(filmsHandler.List))
	mod.RegisterRoute(router.GET, "/{id}", http.HandlerFunc(filmsHandler.Get))

	return mod
}
