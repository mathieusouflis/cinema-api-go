package main

import (
	"authService/db/orm"
	"authService/internal/config"
	"authService/internal/handler"
	pgRepository "authService/internal/repository/postgres"
	redisRepository "authService/internal/repository/redis"
	loginUsecase "authService/internal/usecase/login"
	logoutUsecase "authService/internal/usecase/logout"
	oauthCallbackUsecase "authService/internal/usecase/oauth"
	refreshUsecase "authService/internal/usecase/refresh"
	registerUsecase "authService/internal/usecase/register"
	"context"
	"net/http"
	"os"

	"filmserver/pkg/logger"
	"filmserver/pkg/render"
	"filmserver/pkg/server"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
)

func main() {
	conf := config.Load()
	log := logger.New(conf.Env)

	conn, err := pgx.Connect(context.Background(), conf.Postgres.URL)
	if err != nil {
		log.Error("failed to connect to postgres", "err", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	redisOpts, err := redis.ParseURL(conf.Redis.URL)
	if err != nil {
		log.Error("failed to parse redis URL", "err", err)
		os.Exit(1)
	}
	redisClient := redis.NewClient(redisOpts)
	defer redisClient.Close()

	pg := orm.New(conn)
	userRepo := pgRepository.NewPostgresUserRepository(pg)
	tokenRepo := redisRepository.NewRedisTokenRepository(redisClient)

	deps := handler.Dependencies{
		LoginUseCase:         *loginUsecase.New(userRepo, tokenRepo),
		RegisterUseCase:      *registerUsecase.New(userRepo),
		RefreshUseCase:       *refreshUsecase.New(tokenRepo),
		LogoutUseCase:        *logoutUsecase.New(tokenRepo),
		OauthCallbackUseCase: *oauthCallbackUsecase.New(tokenRepo, userRepo),
	}

	router := chi.NewRouter()
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
