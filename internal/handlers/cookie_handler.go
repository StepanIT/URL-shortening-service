package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/StepanIT/URL-shortening-service/internal/cookies"
	"github.com/gin-gonic/gin"
)

func (h *Handler) SetCookieHandler(c *gin.Context) {
	cookie := http.Cookie{
		Name:     "userID",
		Value:    h.SecretKey,
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	err := cookies.WriteSigned(c.Writer, cookie, []byte(h.SecretKey))
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}

	c.Writer.Write([]byte("cookie set!"))

}

func (h *Handler) GetCookieHandler(c *gin.Context) {
	value, err := cookies.ReadSigned(c.Request, "userID", []byte(h.SecretKey))
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			c.JSON(http.StatusBadRequest, gin.H{"error": "cookie not found"})
		case errors.Is(err, cookies.ErrInvalidValue):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid cookie"})
		default:
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		}
		return
	}

	c.Writer.Write([]byte(value))
}
