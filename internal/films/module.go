package films

import (
	"net/http"

	"example.com/filmserver/pkg/module"
)

func New(server *http.ServeMux) module.Module {
	filmsHandler := NewHandler()

	mod := module.New("/films", server)

	mod.RegisterRoute(module.POST, "/list", http.HandlerFunc(filmsHandler.List))
	mod.RegisterRoute(module.GET, "/{id}", http.HandlerFunc(filmsHandler.Get))

	return mod
}
