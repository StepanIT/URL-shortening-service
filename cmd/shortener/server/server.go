package server

import (
	"net/http"

	"github.com/StepanIT/URL-shortening-service/cmd/shortener/handlers"
	"github.com/StepanIT/URL-shortening-service/cmd/shortener/storage"
)

// функция для создания сервера и обработчиков
func Handler() {
	http.HandleFunc("/", handlers.PostHandler)
	http.HandleFunc("/get/", handlers.GetHandler)

	http.ListenAndServe(storage.ServerAddress, nil)
}
