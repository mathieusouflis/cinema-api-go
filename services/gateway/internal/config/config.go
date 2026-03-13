package config

import "filmserver/pkg/config"

type Config struct {
	config.Base
	config.Postgres
	config.Redis
	config.JWT
	AuthServiceURL string `mapstructure:"AUTH_SERVICE_URL"`
}

func Load() *Config {
	cfg := &Config{}
	if err := config.Load(cfg); err != nil {
		panic("config: " + err.Error())
	}
	config.Required(map[string]string{
		"AUTH_SERVICE_URL": cfg.AuthServiceURL,
	})
	return cfg
}
