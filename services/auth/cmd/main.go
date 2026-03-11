package main

import (
	"authService/db/orm"
	"authService/internal/config"
	"authService/internal/handler"
	repository "authService/internal/repository/postgres"
	usecase "authService/internal/usecase/login"
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

func main() {
	conf := config.Load()

	// Postgres DB connection
	conn, err := pgx.Connect(context.Background(), conf.Postgres.URL)
	if err != nil {
		panic(err)
	}
	defer conn.Close(context.Background())

	pg := orm.New(conn)

	userRepository := repository.NewPostgresUserRepository(pg)

	deps := handler.Dependencies{
		LoginUseCase: *usecase.New(userRepository),
	}

	router := chi.NewRouter()
	handler.Munt(router, &deps)
}
