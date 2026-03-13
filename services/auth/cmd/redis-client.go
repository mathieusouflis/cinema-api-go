package main

import (
	"authService/internal/config"
	"filmserver/pkg/logger"
	"os"

	"github.com/redis/go-redis/v9"
)

func GetRedisClient(conf *config.Config) *redis.Client {
	log := logger.New(conf.Env)
	redisOpts, err := redis.ParseURL(conf.Redis.URL)
	if err != nil {
		log.Error("failed to parse redis URL", "err", err)
		os.Exit(1)
	}
	redisClient := redis.NewClient(redisOpts)
	defer redisClient.Close()
	return redisClient
}
