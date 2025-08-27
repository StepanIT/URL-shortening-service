package handlers

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/StepanIT/URL-shortening-service/internal/cookies"
	"github.com/gin-gonic/gin"
)

type User struct {
	Name string
}

func (h *Handler) SetCookieHandler(c *gin.Context) {
	user := User{Name: "Stepan"}

	var buf bytes.Buffer

	err := gob.NewEncoder(&buf).Encode(&user)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}

	cookie := http.Cookie{
		Name:     "St",
		Value:    buf.String(),
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	err = cookies.WriteSigned(c.Writer, cookie, []byte(h.SecretKey))
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}

	c.Writer.Write([]byte("cookie set!"))

}

func (h *Handler) GetCookieHandler(c *gin.Context) {
	gobEncodeValue, err := cookies.ReadSigned(c.Request, "userID", []byte(h.SecretKey))
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

	var user User

	reader := strings.NewReader(gobEncodeValue)

	if err := gob.NewDecoder(reader).Decode(&user); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}

	fmt.Fprintf(c.Writer, "Name: %q", user.Name)
}
