package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)


type Config struct {
	DBHost string
	DBPort string
	DBUser string
	DBPassword string
	DBName string
	APIUrl string
}


func LoadConfig () (*Config) {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: Could not load .env file: %v\n", err)
	}

	config := &Config{
		DBHost: getEnv("DB_HOST", "localhost"),
		DBPort: getEnv("DB_PORT", "5432"),
		DBUser: getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName: getEnv("DB_NAME", "music"),
		APIUrl: getEnv("API_URL", "http://api-music-info"),
	}

	return config
}

func getEnv (key, defaultValue string) string {
	value, exist := os.LookupEnv(key)
	if !exist {
		return defaultValue
	}
	return value
}