package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) PostShortenHandler(c *gin.Context) {

	// structure for parsing json
	var req struct {
		URL string `json:"url"`
	}

	// parse the request body as JSON and write it to req
	if err := c.BindJSON(&req); err != nil || req.URL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный запрос"})
		return
	}

	// generate a short id and save id+url in storage
	id := generateID()
	if err := h.Repo.Save(id, req.URL); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при сохранении"})
		return
	}

	// create a full address with a short URL
	shortURL := fmt.Sprintf("%s/get/%s", h.BaseURL, id)

	// encode JSON directly via encoding/json
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusCreated, gin.H{"result": shortURL})
}
