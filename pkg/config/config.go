package config

import (
	"time"

	"github.com/go-viper/mapstructure/v2"
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

// Load charge les vars d'env dans la struct fournie.
// Tente de lire .env depuis le répertoire courant (CWD du service).
// Les vraies vars d'env (OS) prennent toujours le dessus sur le fichier .env.
func Load(cfg any) error {
	viper.SetConfigFile(".env")
	viper.SetConfigType("dotenv")
	_ = viper.ReadInConfig() // silencieux si .env absent (prod, CI)

	viper.AutomaticEnv()
	viper.SetDefault("ENV", "development")
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("LOG_LEVEL", "info")

	// viper.AutomaticEnv() only resolves env vars for keys it already knows
	// about. Explicitly bind all known keys so Unmarshal picks them up even
	// when no .env file is present (Docker, CI).
	for _, key := range []string{
		"DATABASE_URL", "DB_MAX_OPEN_CONNS", "DB_MAX_IDLE_CONNS", "DB_CONN_MAX_LIFETIME",
		"REDIS_URL",
		"NATS_URL",
		"JWT_SECRET",
		"ACCESS_TOKEN_TTL", "REFRESH_TOKEN_TTL",
		"S3_ENDPOINT", "S3_REGION", "S3_BUCKET", "S3_ACCESS_KEY", "S3_SECRET_KEY", "S3_USE_SSL",
		"MEILI_URL", "MEILI_MASTER_KEY",
		"TMDB_API_KEY",
		"OAUTH_GOOGLE_CLIENT_ID", "OAUTH_GOOGLE_CLIENT_SECRET",
		"OAUTH_GITHUB_CLIENT_ID", "OAUTH_GITHUB_CLIENT_SECRET",
	} {
		_ = viper.BindEnv(key)
	}

	// Squash=true flattens anonymous embedded structs so that e.g.
	// config.Postgres embedded in a service Config maps DATABASE_URL
	// directly instead of expecting a nested "Postgres" key.
	return viper.Unmarshal(cfg, func(c *mapstructure.DecoderConfig) {
		c.Squash = true
	})
}
