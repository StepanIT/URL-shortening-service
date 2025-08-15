package server

import (
	"log"

	config "github.com/StepanIT/URL-shortening-service"
	"github.com/StepanIT/URL-shortening-service/cmd/shortener/handlers"
	"github.com/StepanIT/URL-shortening-service/cmd/shortener/storage"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// функция для создания сервера и обработчиков
func Handler(cfg *config.Config) {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// interface for working with storage
	var repo storage.Repositories

	// path to file storage
	filePath := cfg.FileStoragePath
	if filePath != "" {
		// use FileStorage
		fs, err := storage.NewFileStorage(filePath)
		if err != nil {
			log.Fatalf("Ошибка создания файлового хранилища: %v", err)
		}
		repo = fs
		log.Println("Используется файловое хранилище:", filePath)
	} else {
		// use in-memory storage
		repo = storage.NewInMemoryStorage()
		log.Println("Используется хранилище в памяти")
	}

	// создаём обработчик, передаём ему выбранное хранилище через интерфейс и конфиг
	h := &handlers.Handler{
		Repo:          repo,
		BaseURL:       cfg.BaseURL,
		ServerAddress: cfg.ServerAddress,
	}

	router := gin.Default()

	router.POST("/", h.PostHandler)
	router.GET("/get/:id", h.GetHandler)

	// created a new route
	router.POST("/api/shorten", h.PostShortenHandler)

	// запуск сервера
	router.Run(cfg.ServerAddress)
}
