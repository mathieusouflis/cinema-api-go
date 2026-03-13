package main

import (
	"filmserver/pkg/logger"
	"filmserver/pkg/server"
	"gateway/internal/config"
	"gateway/internal/handler"
)

func main() {
	conf := config.Load()
	log := logger.New(conf.Env)

	server.Run(":"+conf.Port, handler.NewRouter(&conf.Base), log)
}
