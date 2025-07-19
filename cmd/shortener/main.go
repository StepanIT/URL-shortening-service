package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
)

var urlMap = make(map[int64]string)

func generateId() int64 {
	idRand := rand.Int63()

	fmt.Println("рандомный ID:", idRand)
	return idRand

}

func postHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Читаем тело запроса
		body, err := io.ReadAll(r.Body)
		if err != nil || len(body) == 0 {
			http.Error(w, "ошибка 400", http.StatusBadRequest)
			return
		}
		// преобразуем URL из байтов в строку
		longUrl := string(body)
		// выводим полученный URL
		fmt.Fprintln(w, "Получили URL:", longUrl)
		id := generateId()

		urlMap[id] = longUrl
		fmt.Fprintf(w, "http://localhost:8080/get/%d", id)
	} else {
		http.Error(w, "ошибка 400: не метод POST", http.StatusBadRequest)
	}

}

func getHandler(w http.ResponseWriter, r *http.Request) {

}

func main() {
	generateId()
	http.HandleFunc("/", postHandler)
	http.HandleFunc("/get/", getHandler)

	http.ListenAndServe("localhost:8080", nil)
}
