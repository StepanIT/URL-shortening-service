package main

import (
	"fmt"
	"io"
	"net/http"
)

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
	} else {
		http.Error(w, "ошибка 400: не метод POST", http.StatusBadRequest)
	}

}

func getHandler(w http.ResponseWriter, r *http.Request) {

}

func main() {
	http.HandleFunc("/", postHandler)
	http.HandleFunc("/get", getHandler)

	http.ListenAndServe("localhost:8080", nil)

}
