package server

import (
	"github.com/StepanIT/URL-shortening-service/cmd/shortener/handlers"
	"github.com/StepanIT/URL-shortening-service/cmd/shortener/storage"
	"github.com/gin-gonic/gin"
)

// функция для создания сервера и обработчиков
func Handler() {
	// создаём экхемпляр хранилища
	repo := &storage.InMemoryStorage{
		Data: make(map[string]string),
	}
	// создаём обработчик, передаём ему хранилище через интерфейс
	h := &handlers.Handler{Repo: repo}

	router := gin.Default()

	router.POST("/", h.PostHandler)
	router.GET("/get/:id", h.GetHandler)

	// запуск сервера
	router.Run(storage.ServerAddress)
}
