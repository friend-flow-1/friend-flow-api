package config

import (
	_ "log"
	"os"
	_ "os"

	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBKeyspace string
	JWTSecret  string
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	cfg := &Config{
		DBHost:     GetEnv("DB_HOST", "127.0.0.1:9042"),
		DBKeyspace: GetEnv("DB_KEYSPACE", "friendflow"),
		JWTSecret:  GetEnv("JWT_SECRET", "friendflowsecret"),
	}
	return cfg
}

func GetEnv(key, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return val
}
