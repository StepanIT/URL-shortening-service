package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
)

const serverAdress = "localhost:8080"

// Мапа для хранения ID(сокращенный URL) и полный URL
var urlMap = make(map[string]string)

// функция для генерации ID
func generateID() string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, 7)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "ошибка 400: не метод POST", http.StatusBadRequest)
		return
	}

	// читаем тело запроса
	body, err := io.ReadAll(r.Body)
	if err != nil || len(body) == 0 {
		http.Error(w, "ошибка 400", http.StatusBadRequest)
		return
	}

	// преобразуем URL из байтов в строку
	longURL := string(body)

	// выводим полученный URL
	log.Println("Получили URL:", longURL)

	// получаем ID
	id := generateID()

	// присваеваем полученный URL к полученному ID
	urlMap[id] = longURL

	// выводим ответ с кодом 201 и сокращенный URL
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "http://%s/get/%s", serverAdress, id)

}

func getHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "ошибка 404: не метод GET", http.StatusNotFound)
		return
	}

	// получаем ID из пути запроса, всё что идет после /get/
	id := r.URL.Path[len("/get/"):]
	if id == "" {
		http.Error(w, "ошибка 400: пустой URL", http.StatusBadRequest)
		return
	}

	// ищем оригинальный URL в мапе по полученному ID
	originURL, ok := urlMap[id]
	if !ok {
		http.Error(w, "ошибка 404: URL не найден", http.StatusNotFound)
		return
	}

	// перенаправляет пользователя на оригинальный URL
	http.Redirect(w, r, originURL, http.StatusTemporaryRedirect)

}

// функция для создания сервера и обработчиков
func handler() {
	http.HandleFunc("/", postHandler)
	http.HandleFunc("/get/", getHandler)

	http.ListenAndServe(serverAdress, nil)
}

func main() {
	handler()
}
