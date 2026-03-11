package config

import (
	"time"

	"github.com/spf13/viper"
)

// Base — chaque service embarque ce bloc
type Base struct {
	Env      string `mapstructure:"ENV"` // development | production
	Port     string `mapstructure:"PORT"`
	LogLevel string `mapstructure:"LOG_LEVEL"` // debug | info | warn | error
}

// Postgres — tout service avec une DB
type Postgres struct {
	URL             string        `mapstructure:"DATABASE_URL"`
	MaxOpenConns    int           `mapstructure:"DB_MAX_OPEN_CONNS"`
	MaxIdleConns    int           `mapstructure:"DB_MAX_IDLE_CONNS"`
	ConnMaxLifetime time.Duration `mapstructure:"DB_CONN_MAX_LIFETIME"`
}

func (p *Postgres) Defaults() {
	if p.MaxOpenConns == 0 {
		p.MaxOpenConns = 25
	}
	if p.MaxIdleConns == 0 {
		p.MaxIdleConns = 5
	}
	if p.ConnMaxLifetime == 0 {
		p.ConnMaxLifetime = 5 * time.Minute
	}
}

// Redis
type Redis struct {
	URL string `mapstructure:"REDIS_URL"`
}

// NATS
type NATS struct {
	URL string `mapstructure:"NATS_URL"`
}

// S3 / MinIO
type S3 struct {
	Endpoint  string `mapstructure:"S3_ENDPOINT"`
	Region    string `mapstructure:"S3_REGION"`
	Bucket    string `mapstructure:"S3_BUCKET"`
	AccessKey string `mapstructure:"S3_ACCESS_KEY"`
	SecretKey string `mapstructure:"S3_SECRET_KEY"`
	UseSSL    bool   `mapstructure:"S3_USE_SSL"`
}

// Meilisearch
type Meilisearch struct {
	URL    string `mapstructure:"MEILI_URL"`
	APIKey string `mapstructure:"MEILI_MASTER_KEY"`
}

// JWT — pour la validation dans les services (pas l'émission)
type JWT struct {
	Secret string `mapstructure:"JWT_SECRET"`
}

type TMDB struct {
	APIKey string `mapstructure:"TMDB_API_KEY"`
}

// Load charge les vars d'env dans la struct fournie
func Load(cfg any) error {
	viper.AutomaticEnv()
	viper.SetDefault("ENV", "development")
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("LOG_LEVEL", "info")
	return viper.Unmarshal(cfg)
}
