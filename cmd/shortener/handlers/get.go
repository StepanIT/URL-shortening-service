package handlers

import (
	"net/http"

	"github.com/StepanIT/URL-shortening-service/cmd/shortener/storage"
)

func GetHandler(w http.ResponseWriter, r *http.Request) {
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
	originURL, ok := storage.UrlMap[id]
	if !ok {
		http.Error(w, "ошибка 404: URL не найден", http.StatusNotFound)
		return
	}

	// перенаправляет пользователя на оригинальный URL
	http.Redirect(w, r, originURL, http.StatusTemporaryRedirect)

}
