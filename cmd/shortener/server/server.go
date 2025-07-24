package server

import (
	"net/http"
	"shortener/handlers"
)

const serverAdress = "localhost:8080"

// Мапа для хранения ID(сокращенный URL) и полный URL
var UrlMap = make(map[string]string)

// функция для создания сервера и обработчиков
func Handler() {
	http.HandleFunc("/", handlers.PostHandler)
	http.HandleFunc("/get/", handlers.GetHandler)

	http.ListenAndServe(serverAdress, nil)
}