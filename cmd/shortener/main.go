package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
)

// Мапа для хранения ID(сокращенный URL) и полный URL
var urlMap = make(map[string]string)

// функция для генерации ID
func generateId() string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, 7)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {

		// читаем тело запроса
		body, err := io.ReadAll(r.Body)
		if err != nil || len(body) == 0 {
			http.Error(w, "ошибка 400", http.StatusBadRequest)
			return
		}

		// преобразуем URL из байтов в строку
		longUrl := string(body)

		// выводим полученный URL
		fmt.Fprintln(w, "Получили URL:", longUrl)

		// получаем ID
		id := generateId()

		// присваеваем полученный URL к полученному ID
		urlMap[id] = longUrl

		// выводим ответ с кодом 201 и сокращенный URL
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "http://localhost:8080/get/%s\n", id)

		fmt.Fprintln(w, "мапа", urlMap)
	} else {
		http.Error(w, "ошибка 400: не метод POST", http.StatusBadRequest)
	}

}

func getHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {

		// получаем ID из пути запроса, всё что идет после /get/
		id := r.URL.Path[len("/get/"):]
		if id == "" {
			http.Error(w, "ошибка 400: пустой URL", http.StatusBadRequest)
			return
		}

		// ищем оригинальный URL в мапе по полученному ID
		originUrl, ok := urlMap[id]
		if !ok {
			http.Error(w, "ошибка 404: URL не найден", http.StatusNotFound)
			return
		}

		// перенаправляет пользователя на оригинальный URL
		http.Redirect(w, r, originUrl, http.StatusTemporaryRedirect)
	} else {
		http.Error(w, "ошибка 400: не метод GET", http.StatusBadRequest)
	}
}

// функция для создания сервера и обработчиков
func handler() {
	http.HandleFunc("/", postHandler)
	http.HandleFunc("/get/", getHandler)

	http.ListenAndServe("localhost:8080", nil)
}

func main() {
	handler()
}
