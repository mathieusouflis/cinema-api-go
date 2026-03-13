package main

import (
	"authService/internal/config"
	"authService/internal/handler"
	"net/http"

	"filmserver/pkg/logger"
	"filmserver/pkg/middleware"
	"filmserver/pkg/render"
	"filmserver/pkg/server"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func main() {
	conf := config.Load()
	log := logger.New(conf.Env)

	pg := GetPgClient(conf)
	redisClient := GetRedisClient(conf)

	deps := GetDependencies(conf, pg, redisClient)

	router := chi.NewRouter()
	router.Use(middleware.NewLogMiddleware())
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type"},
		AllowCredentials: true,
	}))
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, http.StatusOK, map[string]string{"message": "Hello, World!"})
	})
	handler.Munt(router, &deps)
	server.Run(":"+conf.Port, router, log)
}
