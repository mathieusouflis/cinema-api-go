package server

import (
	"net/http"

	"example.com/filmserver/internal/auth"
	"example.com/filmserver/internal/films"
)

func New() http.Handler {
	mux := http.NewServeMux()

	authHandler := auth.NewHandler()
	filmsHandler := films.NewHandler()

	mux.HandleFunc("POST /auth/login", authHandler.Login)
	mux.HandleFunc("POST /auth/register", authHandler.Register)

	mux.Handle("GET /films", auth.Middleware(http.HandlerFunc(filmsHandler.List)))
	mux.Handle("GET /films/{id}", auth.Middleware(http.HandlerFunc(filmsHandler.Get)))

	return mux
}
