package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type BatchRequestItem struct {
	CorrelationID string `json:"correlation_id" binding:"required"`
	OriginalURL   string `json:"original_url" binding:"required,url"`
}

type BatchResponseItem struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

func (h *Handler) PostShortenBatchHandler(c *gin.Context) {
	var reqItems []BatchRequestItem

	if err := c.ShouldBindJSON(&reqItems); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request format"})
	}

	userID := c.GetString("userID")

	respItems := make([]BatchResponseItem, 0, len(reqItems))

	for _, item := range reqItems {
		id := generateID()
		err := h.Repo.Save(id, item.OriginalURL, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to shorten URL"})
			return
		}

		respItems = append(respItems, BatchResponseItem{
			CorrelationID: item.CorrelationID,
			ShortURL:      h.BaseURL + "/" + id,
		})
	}
	c.JSON(http.StatusCreated, respItems)

}
