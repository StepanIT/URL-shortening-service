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
	var batchRequests []BatchRequestItem

	if err := c.ShouldBindJSON(&batchRequests); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request format"})
		return
	}

	userID, exists := c.Get("userID")
	batchResponses := make([]BatchResponseItem, 0, len(batchRequests))

	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user not identified"})
		return
	}

	for _, item := range batchRequests {
		id := generateID()
		err := h.Repo.Save(id, item.OriginalURL, userID.(string))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to shorten URL"})
			return
		}

		batchResponses = append(batchResponses, BatchResponseItem{
			CorrelationID: item.CorrelationID,
			ShortURL:      h.BaseURL + "/" + id,
		})
	}
	c.JSON(http.StatusCreated, batchResponses)

}
