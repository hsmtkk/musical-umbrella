package env

import (
	"log"
	"os"
)

func RequiredEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("environment variable %s must be defined", key)
	}
	return val
}
