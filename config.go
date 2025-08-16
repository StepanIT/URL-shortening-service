package config

import (
	"flag"
	"log"
	"os"

	"github.com/joho/godotenv"
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

	// default values
	defaultServerAddress := ":8080"
	defaultBaseURL := "http://localhost:8080"
	defaultFileStoragePath := ""

	// define flags
	flagServerAddress := flag.String("server-port", defaultServerAddress, "адрес запуска HTTP-сервера")
	flagBaseURL := flag.String("b", defaultBaseURL, " базовый адрес результирующего сокращённого URL")
	flagFileStoragePath := flag.String("file-storage-path", defaultFileStoragePath, "путь до файла с сокращёнными URL")
	flag.Parse()

	// get value from ENV or use flag
	serverAddress := "localhost" + getEnvOrFlag("SERVER_PORT", *flagServerAddress, defaultFileStoragePath)
	baseURL := getEnvOrFlag("BASE_URL", *flagBaseURL, defaultBaseURL)
	fileStoragePath := getEnvOrFlag("FILE_STORAGE_PATH", *flagFileStoragePath, defaultFileStoragePath)

	return &Config{
		ServerAddress:   serverAddress,
		BaseURL:         baseURL,
		FileStoragePath: fileStoragePath,
	}
}

// returns the value of the environment variable key
func getEnvOrFlag(envKey string, flagValue string, defaultVal string) string {
	if value, exist := os.LookupEnv(envKey); exist && value != "" {
		return value
	}
	if flagValue != "" {
		return flagValue
	}

	return defaultVal
}
