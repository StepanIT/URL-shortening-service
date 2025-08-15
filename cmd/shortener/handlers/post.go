package handlers

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"

	"github.com/StepanIT/URL-shortening-service/cmd/shortener/storage"
	"github.com/gin-gonic/gin"
)

// структура с интерфейсом для работы с хранилищем
type Handler struct {
	Repo          storage.Repositories
	BaseURL       string
	ServerAddress string
}

// функция для генерации ID
func generateID() string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, 7)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func (h *Handler) PostHandler(c *gin.Context) {

	// читаем тело запроса
	body, err := io.ReadAll(c.Request.Body)
	if err != nil || len(body) == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "ошибка 400"})
		return
	}
	log.Println("Получили URL:", body)

	// преобразуем URL из байтов в строку
	LongURL := string(body)

	// выводим полученный URL
	log.Println("Получили URL:", LongURL)

	// получаем ID
	id := generateID()

	// присваеваем полученный URL к полученному ID
	err = h.Repo.Save(id, LongURL)
	log.Println("Присвоенный URL", err)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при сохранении"})
		return
	}

	// сhecking that BaseURL is installed
	base := h.BaseURL
	if base == "" {
		base = "http://localhost:8080"
	}

	// выводим ответ с кодом 201 и сокращенный URL
	shortURL := fmt.Sprintf("%s/get/%s", base, id)
	c.JSON(http.StatusCreated, gin.H{"result": shortURL})

}
