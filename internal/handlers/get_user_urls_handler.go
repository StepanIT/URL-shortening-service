package handlers

import (
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

	// Формируем ответ в нужном формате
	response := make([]struct {
		ShortURL    string `json:"short_url"`
		OriginalURL string `json:"original_url"`
	}, 0, len(userURLs))

	for shortID, originalURL := range userURLs {
		response = append(response, struct {
			ShortURL    string `json:"short_url"`
			OriginalURL string `json:"original_url"`
		}{
			ShortURL:    h.BaseURL + "/get/" + shortID,
			OriginalURL: originalURL,
		})
	}

	// Устанавливаем правильный Content-Type и возвращаем JSON
	c.JSON(http.StatusOK, response)
}
