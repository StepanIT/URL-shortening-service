package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// structure for storing server settings
type Config struct {
	ServerAddress string
	BaseURL       string
}

// create a configuration by loading environment variables from .env
func NewConfig() *Config {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	return &Config{
		ServerAddress: getEnv("SERVER_ADDRESS", "localhost:8080"),
		BaseURL:       getEnv("BASE_URL", "http://localhost:8080"),
	}
}

// returns the value of the environment variable key
func getEnv(key string, defaultVal string) string {
	if value, exist := os.LookupEnv(key); exist {
		return value
	}

	return defaultVal
}
