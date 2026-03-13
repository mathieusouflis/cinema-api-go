package handler

import (
	"filmserver/pkg/config"
	"filmserver/pkg/middleware"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewRouter(cfg *config.Base) *chi.Mux {
	router := chi.NewRouter()

	//RATE LIMIT MIDDLEWARE
	//LOG MIDDLEWARE
	router.Use(middleware.NewLogMiddleware())

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "ok")
	})

	NewAuthHandler(router)

	//GRAPHANA MIDDLEWARE

	return router
}
