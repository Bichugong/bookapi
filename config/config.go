package config

import (
	"log"
	"os"
	"github.com/joho/godotenv"
)

type Config struct {
	JWTSecret string
	DBURL     string
}

func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found")
	}

	return Config{
		JWTSecret: getEnv("JWT_SECRET", "default_secret_key"),
		DBURL:     getEnv("DB_URL", ""),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}