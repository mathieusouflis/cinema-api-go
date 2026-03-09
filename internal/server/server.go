package server

import (
	"net/http"

	"example.com/filmserver/internal/auth"
	"example.com/filmserver/internal/films"
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
