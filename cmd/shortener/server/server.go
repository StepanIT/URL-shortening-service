package server

import (
	"net/http"

	"github.com/StepanIT/URL-shortening-service/cmd/shortener/handlers"
	"github.com/StepanIT/URL-shortening-service/cmd/shortener/storage"
)

// функция для создания сервера и обработчиков
func Handler() {
	// создаём экхемпляр хранилища
	repo := &storage.InMemoryStorage{
		Data: make(map[string]string),
	}
	// создаём обработчик, передаём ему хранилище через интерфейс
	h := &handlers.Handler{Repo: repo}

	http.HandleFunc("/", h.PostHandler)
	http.HandleFunc("/get/", h.GetHandler)

	// запуск сервера
	http.ListenAndServe(storage.ServerAddress, nil)
}
