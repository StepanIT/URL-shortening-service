package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetUserURLsHandler возвращает все URL пользователя
func (h *Handler) GetUserURLsHandler(c *gin.Context) {
	userIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	userID := userIDValue.(string)

	urlsMap, err := h.Repo.GetAllByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch URLs"})
		return
	}

	if len(urlsMap) == 0 {
		c.Status(http.StatusNoContent)
		return
	}

	// Преобразуем map в срез объектов для ответа
	type respPair struct {
		ShortURL    string `json:"short_url"`
		OriginalURL string `json:"original_url"`
	}

	resp := make([]respPair, 0, len(urlsMap))
	for id, original := range urlsMap {
		resp = append(resp, respPair{
			ShortURL:    h.BaseURL + "/get/" + id,
			OriginalURL: original,
		})
	}

	c.JSON(http.StatusOK, resp)
}
