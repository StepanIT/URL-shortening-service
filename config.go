package config

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerAddress   string
	BaseURL         string
	FileStoragePath string
}

func NewConfig() *Config {
	// Загрузка .env файла (игнорируем ошибку, если файла нет)
	_ = godotenv.Load()

	// Установка значений по умолчанию
	defaults := map[string]string{
		"SERVER_ADDRESS":    "localhost:8080",
		"BASE_URL":          "http://localhost:8080",
		"FILE_STORAGE_PATH": "",
	}

	// Определение флагов командной строки
	flags := map[string]*string{
		"a": flag.String("a", defaults["SERVER_ADDRESS"], "HTTP server address"),
		"b": flag.String("b", defaults["BASE_URL"], "Base URL for shortened links"),
		"f": flag.String("f", defaults["FILE_STORAGE_PATH"], "Path to storage file"),
	}
	flag.Parse()

	// Получение финальных значений (env > flags > defaults)
	cfg := &Config{
		ServerAddress:   getValue("SERVER_ADDRESS", flags["a"], defaults["SERVER_ADDRESS"]),
		BaseURL:         getValue("BASE_URL", flags["b"], defaults["BASE_URL"]),
		FileStoragePath: getValue("FILE_STORAGE_PATH", flags["f"], defaults["FILE_STORAGE_PATH"]),
	}

	// Нормализация адреса сервера (добавляем : если нужно)
	if !strings.Contains(cfg.ServerAddress, ":") {
		cfg.ServerAddress = ":" + cfg.ServerAddress
	}

	log.Printf("Config loaded: %+v", cfg)
	return cfg
}

func getValue(envKey string, flagValue *string, defaultValue string) string {
	// Приоритет: ENV > Flag > Default
	if val, exists := os.LookupEnv(envKey); exists {
		return val
	}
	if flagValue != nil && *flagValue != "" {
		return *flagValue
	}
	return defaultValue
}
