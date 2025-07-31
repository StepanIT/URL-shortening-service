package handlers

import (
	"log"
	"net/http"
)

func (h *Handler) GetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "ошибка 404: не метод GET", http.StatusNotFound)
		return
	}

	// получаем ID из пути запроса, всё что идет после /get/
	id := r.URL.Path[len("/get/"):]
	log.Println("найденный ID", id)
	if id == "" {
		http.Error(w, "ошибка 400: пустой URL", http.StatusBadRequest)
		return
	}

	// ищем оригинальный URL по ID через метод Get интерфейса Repo
	LongURL, err := h.Repo.Get(id)
	log.Println("найденный URL", LongURL)
	if err != nil {
		http.Error(w, "ошибка 404: URL не найден", http.StatusNotFound)
		return
	}

	// перенаправляет пользователя на оригинальный URL
	http.Redirect(w, r, LongURL, http.StatusTemporaryRedirect)

}
