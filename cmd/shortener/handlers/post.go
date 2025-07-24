package handlers

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"

	"github.com/StepanIT/URL-shortening-service/cmd/shortener/storage"
)

// функция для генерации ID
func generateID() string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, 7)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
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
	storage.UrlMap[id] = longURL

	// выводим ответ с кодом 201 и сокращенный URL
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "http://%s/get/%s", storage.ServerAddress, id)

}
