package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetUserURLsHandler(c *gin.Context) {
	// Получаем userID из контекста
	userIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user not identified"})
		return
	}

	userID := userIDValue.(string)

	// Используем новый метод хранилища
	userURLs, err := h.Repo.GetAllByUser(userID)
	if err != nil {
		log.Printf("Error getting user URLs: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	// Если у пользователя нет ссылок, возвращаем 204
	if len(userURLs) == 0 {
		c.Status(http.StatusNoContent)
		return
	}

	// Формируем ответ в нужном формате
	var response []struct {
		ShortURL    string `json:"short_url"`
		OriginalURL string `json:"original_url"`
	}

	for shortID, originalURL := range userURLs {
		response = append(response, struct {
			ShortURL    string `json:"short_url"`
			OriginalURL string `json:"original_url"`
		}{
			ShortURL:    fmt.Sprintf("%s/%s", h.BaseURL, shortID),
			OriginalURL: originalURL,
		})
	}

	// Устанавливаем правильный Content-Type и возвращаем JSON
	c.JSON(http.StatusOK, response)
}
