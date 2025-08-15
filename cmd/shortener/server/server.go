package server

import (
	"log"

	config "github.com/StepanIT/URL-shortening-service"
	"github.com/StepanIT/URL-shortening-service/cmd/shortener/handlers"
	"github.com/StepanIT/URL-shortening-service/cmd/shortener/storage"
	"github.com/gin-gonic/gin"
)

func Handler(cfg *config.Config) {
	// Инициализация хранилища
	repo, err := initStorage(cfg.FileStoragePath)
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}

	// Создание обработчика
	h := &handlers.Handler{
		Repo:    repo,
		BaseURL: cfg.BaseURL,
	}

	// Настройка роутера
	router := setupRouter(h)

	// Запуск сервера с логированием
	log.Printf("Starting server on %s", cfg.ServerAddress)
	log.Printf("Base URL: %s", cfg.BaseURL)
	if cfg.FileStoragePath != "" {
		log.Printf("Using file storage: %s", cfg.FileStoragePath)
	} else {
		log.Println("Using in-memory storage")
	}

	if err := router.Run(cfg.ServerAddress); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func initStorage(filePath string) (storage.Repositories, error) {
	if filePath != "" {
		fs, err := storage.NewFileStorage(filePath)
		if err != nil {
			return nil, err
		}
		return fs, nil
	}
	return storage.NewInMemoryStorage(), nil
}

func setupRouter(h *handlers.Handler) *gin.Engine {
	router := gin.Default()

	// API v1
	v1 := router.Group("/")
	{
		v1.POST("/", h.PostHandler)                   // TEXT → URL
		v1.GET("/get/:id", h.GetHandler)              // Редирект
		v1.POST("/api/shorten", h.PostShortenHandler) // JSON → URL
	}

	return router
}
