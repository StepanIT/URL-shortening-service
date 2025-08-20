package handlers

import (
	"encoding/json"
	"fmt"
	"log"
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

	log.Println("base URL", h.BaseURL)
	shortURL := fmt.Sprintf("%s/get/%s", h.BaseURL, id)

	resp := struct {
		Result string `json:"result"`
	}{Result: shortURL}

	c.Header("Content-Type", "application/json")
	c.Writer.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(c.Writer).Encode(resp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при кодировании JSON"})
	}
}
