package server

import (
	"flag"
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

	addr := flag.String("a", cfg.ServerAddress, "server address")
	base := flag.String("b", cfg.BaseURL, "base URL")
	file := flag.String("f", cfg.FileStoragePath, "file storage path")
	flag.Parse()

	// interface for working with storage
	var repo storage.Repositories
	if *file != "" {
		fs, err := storage.NewFileStorage(*file)
		if err != nil {
			log.Fatalf("Ошибка создания файлового хранилища: %v", err)
		}
		repo = fs
		log.Println("Используется файловое хранилище:", *file)
	} else {
		repo = storage.NewInMemoryStorage()
		log.Println("Используется хранилище в памяти")
	}

	// создаём обработчик, передаём ему выбранное хранилище через интерфейс и конфиг
	h := &handlers.Handler{
		Repo:          repo,
		BaseURL:       *base,
		ServerAddress: *addr,
	}
	router := gin.Default()

	router.POST("/", h.PostHandler)
	router.GET("/get/:id", h.GetHandler)

	// created a new route
	router.POST("/api/shorten", h.PostShortenHandler)

	// запуск сервера
	router.Run(*addr)
}
