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
func Handler() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	cfg := config.NewConfig()

	// создаём экхемпляр хранилища
	repo := storage.NewInMemoryStorage()
	// создаём обработчик, передаём ему хранилище через интерфейс и конфиг
	h := &handlers.Handler{
		Repo:    repo,
		BaseURL: cfg.BaseURL,
	}

	router := gin.Default()

	router.POST("/", h.PostHandler)
	router.GET("/get/:id", h.GetHandler)

	// created a new route
	router.POST("/api/shorten", h.PostShortenHandler)

	// запуск сервера
	router.Run(cfg.ServerAddress)
}
