package middleware

import (
	"log"
	"math/rand"
	"net/http"

	"github.com/StepanIT/URL-shortening-service/internal/cookies"
	"github.com/gin-gonic/gin"
)

type User struct {
	ID string
}

// generateID генерирует случайный идентификатор (аналогично функции в handlers)
func generateUserID() string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 16)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// Auth middleware проверяет и устанавливает куку пользователя
func Auth(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var userID string
		var needSetCookie bool

		// Пытаемся прочитать подписанную куку
		cookieValue, err := cookies.ReadSigned(c.Request, "userID", secretKey)
		if err != nil || cookieValue == "" {
			// Если нет куки или подпись невалидна, создаём новый userID
			userID = generateUserID()
			needSetCookie = true
			log.Println("Cookie invalid or not found, generating new user ID:", userID)
		} else {
			userID = cookieValue
			log.Println("User authenticated:", userID)
		}

		// Сохраняем userID в контексте Gin
		c.Set("userID", userID)

		// Если куки не было или она была невалидна, устанавливаем новую
		if needSetCookie {

			cookie := http.Cookie{
				Name:     "userID",
				Value:    userID,
				Path:     "/",
				MaxAge:   3600 * 24 * 30,
				HttpOnly: true,
				Secure:   false,
				SameSite: http.SameSiteLaxMode,
			}

			// Подписываем и устанавливаем куку
			if err := cookies.WriteSigned(c.Writer, cookie, secretKey); err != nil {
				log.Printf("Failed to set signed cookie: %v", err)
			} else {
				log.Printf("New cookie set for user: %s", userID)
			}
		}

		c.Next()
	}
}
