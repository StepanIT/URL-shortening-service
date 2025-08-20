package config

import (
	"flag"
	"log"
	"os"

	"github.com/joho/godotenv"
)

const (
	defaultServerAddress   = "localhost:8080"
	defaultBaseURL         = "http://localhost:8080"
	defaultFileStoragePath = ""
)

// structure for storing server settings
type Config struct {
	ServerAddress   string
	BaseURL         string
	FileStoragePath string
}

// create a configuration by loading environment variables from .env
func NewConfig() *Config {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	// define flags
	flagServerAddress := flag.String("a", defaultServerAddress, "адрес запуска HTTP-сервера")
	flagBaseURL := flag.String("b", defaultBaseURL, " базовый адрес результирующего сокращённого URL")
	flagFileStoragePath := flag.String("f", defaultFileStoragePath, "путь до файла с сокращёнными URL")
	flag.Parse()

	return &Config{
		ServerAddress:   getEnvOrFlag("SERVER_ADDRESS", flagServerAddress, defaultServerAddress),
		BaseURL:         getEnvOrFlag("BASE_URL", flagBaseURL, defaultBaseURL),
		FileStoragePath: getEnvOrFlag("FILE_STORAGE_PATH", flagFileStoragePath, defaultFileStoragePath),
	}
}

// returns the value of the environment variable key
func getEnvOrFlag(envKey string, flagValue *string, defaultVal string) string {
	if value, exists := os.LookupEnv(envKey); exists && value != "" {
		return value
	}
	if flagValue != nil && *flagValue != defaultVal {
		return *flagValue
	}
	return defaultVal
}
