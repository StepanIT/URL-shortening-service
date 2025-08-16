package config

import (
	"flag"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config хранит настройки сервера
type Config struct {
	ServerAddress   string // host:port для запуска сервера
	BaseURL         string // базовый URL для сокращённых ссылок
	FileStoragePath string // путь к файлу с сохранёнными URL
}

// NewConfig создаёт конфигурацию сервера, учитывая ENV и флаги
func NewConfig() *Config {
	// загрузка из .env (если есть)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// default значения
	defaultHost := "localhost"
	defaultPort := "8080"
	defaultBaseURL := ""
	defaultFileStoragePath := ""

	// флаги
	flagServerHost := flag.String("server-host", defaultHost, "адрес хоста для запуска HTTP-сервера")
	flagServerPort := flag.String("server-port", defaultPort, "порт для запуска HTTP-сервера")
	flagBaseURL := flag.String("server-base-url", defaultBaseURL, "базовый адрес для сокращённых URL")
	flagFileStoragePath := flag.String("file-storage-path", defaultFileStoragePath, "путь до файла с сокращёнными URL")
	flag.Parse()

	// ENV или флаги
	host := getEnvOrFlag("SERVER_HOST", *flagServerHost, defaultHost)
	port := getEnvOrFlag("SERVER_PORT", *flagServerPort, defaultPort)
	serverAddress := host + ":" + port

	baseURL := getEnvOrFlag("BASE_URL", *flagBaseURL, "")
	if baseURL == "" {
		baseURL = "http://" + serverAddress
	}

	fileStoragePath := getEnvOrFlag("FILE_STORAGE_PATH", *flagFileStoragePath, defaultFileStoragePath)

	return &Config{
		ServerAddress:   serverAddress,
		BaseURL:         baseURL,
		FileStoragePath: fileStoragePath,
	}
}

// getEnvOrFlag отдаёт значение из ENV или флага, иначе default
func getEnvOrFlag(envKey, flagValue, defaultVal string) string {
	if value, exist := os.LookupEnv(envKey); exist && value != "" {
		return value
	}
	if flagValue != "" {
		return flagValue
	}
	return defaultVal
}
