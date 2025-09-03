package middleware

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"github.com/StepanIT/URL-shortening-service/internal/cookies"
	"github.com/gin-gonic/gin"
)

type User struct {
	ID string
}

func Auth(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := cookies.ReadSigned(c.Request, "userID", []byte(secretKey))
		if err != nil || userID == "" {
			// Если нет куки или она недействительна, генерируем новую
			userID = generateRandomID()
			cookie := http.Cookie{
				Name:     "userID",
				Value:    userID,
				Path:     "/",
				HttpOnly: true,
			}
			_ = cookies.WriteSigned(c.Writer, cookie, []byte(secretKey))
		}

		// Сохраняем userID в контексте GIN
		c.Set("userID", userID)
		c.Next()
	}
}

// Генерация случайного ID для пользователя
func generateRandomID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}
