package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetUserURLsHandler(c *gin.Context) {
	// Получаем userID из контекста, установленного middleware
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user not identified"})
		return
	}

	// Используем новый метод хранилища
	userURLs, err := h.Repo.GetAllByUser(userID.(string))
	if err != nil {
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
			ShortURL:    fmt.Sprintf("%s/get/%s", h.BaseURL, shortID),
			OriginalURL: originalURL,
		})
	}

	c.JSON(http.StatusOK, response)
}
