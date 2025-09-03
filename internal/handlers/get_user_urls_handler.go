package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetUserURLsHandler(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "userID not found"})
		return
	}

	urls, _ := h.Repo.GetAllByUser(userID.(string))
	if len(urls) == 0 {
		c.Status(http.StatusNoContent)
		return
	}

	resp := make([]respPair, 0, len(urls))
	for id, original := range urls {
		resp = append(resp, respPair{
			ShortURL:    h.BaseURL + "/get/" + id,
			OriginalURL: original,
		})
	}

	c.JSON(http.StatusOK, resp)
}
