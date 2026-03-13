package middleware

import (
	"filmserver/pkg/config"
	"filmserver/pkg/logger"
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	status int
}

type Config struct {
	Base config.Base
}

func NewLogMiddleware() func(next http.Handler) http.Handler {
	cfg := &Config{}
	err := config.Load(cfg)

	if err != nil {
		panic(err)
	}

	log := logger.New(cfg.Base.Env)

	log.Info("log middleware loaded")
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}
			start := time.Now()
			next.ServeHTTP(rw, r)
			log.Info(r.Method + " " + r.URL.Path + " - " + http.StatusText(rw.status) + " " + time.Since(start).String())
		})
	}
}
