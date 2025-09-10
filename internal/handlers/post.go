package handlers

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"

	"github.com/StepanIT/URL-shortening-service/internal/storage"
	"github.com/gin-gonic/gin"
)

// структура с интерфейсом для работы с хранилищем
type Handler struct {
	Repo          storage.URLShortenerRepositories
	BaseURL       string
	ServerAddress string
	SecretKey     string
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
	body, err := io.ReadAll(c.Request.Body)
	if err != nil || len(body) == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "ошибка 400"})
		return
	}
	LongURL := string(body)

	userIDValue, exists := c.Get("userID")
	if !exists {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "ошибка аутентификации"})
		return
	}
	userID := userIDValue.(string)

	id := generateID()
	shortURL, err := h.Repo.Save(id, LongURL, userID)
	if err != nil {
		if err.Error() == "url already exists" {
			c.String(http.StatusConflict, fmt.Sprintf("%s/get/%s", h.BaseURL, shortURL))
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при сохранении"})
		return
	}

	c.String(http.StatusCreated, fmt.Sprintf("%s/get/%s", h.BaseURL, shortURL))
}
