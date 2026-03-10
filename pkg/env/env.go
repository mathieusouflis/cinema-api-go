package env

import "os"

func GetEnv(env string, fallback string) string {
	value := os.Getenv(env)
	if value == "" {
		return fallback
	}
	return value
}
