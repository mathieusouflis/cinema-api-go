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

	"github.com/redis/go-redis/v9"
)

func GetDependencies(conf *config.Config, pgClient *orm.Queries, redisClient *redis.Client) handler.Dependencies {

	userRepo := pgRepository.NewPostgresUserRepository(pgClient)
	tokenRepo := redisRepository.NewRedisTokenRepository(redisClient)
	deps := handler.Dependencies{
		LoginUseCase:         *loginUsecase.New(userRepo, tokenRepo),
		RegisterUseCase:      *registerUsecase.New(userRepo),
		RefreshUseCase:       *refreshUsecase.New(tokenRepo),
		LogoutUseCase:        *logoutUsecase.New(tokenRepo),
		OauthCallbackUseCase: *oauthCallbackUsecase.New(tokenRepo, userRepo),
	}

	return deps
}
