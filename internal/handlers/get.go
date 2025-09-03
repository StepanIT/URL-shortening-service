package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetHandler(c *gin.Context) {
	id := c.Param("id")
	url, err := h.Repo.Get(id)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	c.Redirect(http.StatusTemporaryRedirect, url)
}
