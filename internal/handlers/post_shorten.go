package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) PostShortenHandler(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "userID not found"})
		return
	}

	var req struct {
		URL string `json:"url"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	shortID := generateID() // твоя функция генерации короткого ID
	if err := h.Repo.Save(shortID, req.URL, userID.(string)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot save URL"})
		return
	}

	c.JSON(http.StatusCreated, respPair{
		ShortURL:    h.BaseURL + "/get/" + shortID,
		OriginalURL: req.URL,
	})
}
