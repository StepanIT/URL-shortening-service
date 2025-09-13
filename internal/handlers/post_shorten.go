package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) PostShortenHandler(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user not identified"})
		return
	}

	var req struct {
		URL string `json:"url"`
	}
	if err := c.BindJSON(&req); err != nil || req.URL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный запрос"})
		return
	}

	id := generateID()
	shortURL, err := h.Repo.Save(id, req.URL, userID.(string))
	if err != nil {
		if err.Error() == "url already exists" {
			// возвращаем существующий URL с кодом 409
			resp := struct {
				Result string `json:"result"`
			}{Result: fmt.Sprintf("%s/get/%s", h.BaseURL, shortURL)}
			c.JSON(http.StatusConflict, resp)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при сохранении"})
		return
	}

	resp := struct {
		Result string `json:"result"`
	}{Result: fmt.Sprintf("%s/get/%s", h.BaseURL, shortURL)}

	c.JSON(http.StatusCreated, resp)
}
