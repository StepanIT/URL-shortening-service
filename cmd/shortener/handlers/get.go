package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetHandler(c *gin.Context) {
	// получаем ID из пути запроса, всё что идет после /get/
	id := c.Param("id")
	if id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "ошибка 400: пустой URL"})
		return
	}
	// ищем оригинальный URL по ID через метод Get интерфейса Repo
	LongURL, err := h.Repo.Get(id)
	log.Println("найденный URL", LongURL)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": "ошибка 404: URL не найден"})
		return
	}

	// перенаправляет пользователя на оригинальный URL
	c.Redirect(http.StatusTemporaryRedirect, LongURL)

}
