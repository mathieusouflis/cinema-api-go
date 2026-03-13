package main

import (
	"authService/db/orm"
	"authService/internal/config"
	"context"
	"filmserver/pkg/logger"
	"os"

	"github.com/jackc/pgx/v5"
)

func GetPgClient(conf *config.Config) *orm.Queries {
	log := logger.New(conf.Env)
	conn, err := pgx.Connect(context.Background(), conf.Postgres.URL)
	if err != nil {
		log.Error("failed to connect to postgres", "err", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())
	pg := orm.New(conn)
	return pg
}
