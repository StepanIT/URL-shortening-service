package server

import (
	"log"

	"github.com/StepanIT/URL-shortening-service/internal/config"
	"github.com/StepanIT/URL-shortening-service/internal/handlers"
	"github.com/StepanIT/URL-shortening-service/internal/storage"
	"github.com/gin-gonic/gin"
)

// function for creating server and handlers
func Handler(cfg *config.Config) {
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

	// pass the selected storage and config to the handler
	h := &handlers.Handler{
		Repo:          repo,
		BaseURL:       cfg.BaseURL,
		ServerAddress: cfg.ServerAddress,
	}

	// setting up GIN routes
	router := gin.Default()

	router.POST("/", h.PostHandler)
	router.GET("/get/:id", h.GetHandler)
	router.POST("/api/shorten", h.PostShortenHandler)

	// starting the server
	router.Run(cfg.ServerAddress)
}
