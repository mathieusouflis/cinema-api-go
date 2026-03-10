package server

import (
	"net/http"

	"example.com/filmserver/api/auth"
	"example.com/filmserver/api/films"
)

func New() http.Handler {
	mux := http.NewServeMux()

	authModule := auth.New(mux)
	filmsModule := films.New(mux)

	authModule.Register()
	filmsModule.Register()

	filmsModule.PrintRoutesDocumentation()
	authModule.PrintRoutesDocumentation()

	return mux
}
