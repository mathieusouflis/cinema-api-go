package config

import (
	"filmserver/pkg/config"
)

type Config struct {
	config.Base
	config.Postgres
	config.Redis
	config.NATS
	config.JWT
}

func Load() *Config {
	cfg := &Config{}
	if err := config.Load(cfg); err != nil {
		panic("config: " + err.Error())
	}
	cfg.Postgres.Defaults()
	config.Required(map[string]string{
		"DATABASE_URL": cfg.Postgres.URL,
		"REDIS_URL":    cfg.Redis.URL,
		"JWT_SECRET":   cfg.JWT.Secret,
		"NATS_URL":     cfg.NATS.URL,
	})
	return cfg
}
