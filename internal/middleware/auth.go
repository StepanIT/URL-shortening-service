package middleware

import (
	"bytes"
	"encoding/gob"
	"log"
	"math/rand"
	"net/http"
	"strings"

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
		var user User
		var userID string
		var needSetCookie bool = false

		// Пытаемся прочитать и верифицировать куку
		cookieValue, err := cookies.ReadSigned(c.Request, "userID", []byte(secretKey))
		if err != nil {
			// Если куки нет или она невалидна, создаем нового пользователя
			log.Println("Cookie invalid or not found, generating new user ID")
			userID = generateUserID()
			user = User{ID: userID}
			needSetCookie = true
		} else {
			// Если кука есть, пробуем декодировать из Gob
			reader := strings.NewReader(cookieValue)
			if err := gob.NewDecoder(reader).Decode(&user); err != nil {
				log.Printf("Failed to decode user from cookie: %v", err)
				// Если не удалось декодировать, создаем нового пользователя
				userID = generateUserID()
				user = User{ID: userID}
				needSetCookie = true
			} else {
				userID = user.ID
				log.Printf("User authenticated: %s", userID)
			}
		}

		// Сохраняем userID в контексте Gin, чтобы хендлеры могли его использовать
		c.Set("userID", userID)

		// Если куки не было или она была невалидна, устанавливаем новую
		if needSetCookie {
			// Кодируем структуру User в Gob
			var buf bytes.Buffer
			err := gob.NewEncoder(&buf).Encode(&user)
			if err != nil {
				log.Printf("Failed to encode user for cookie: %v", err)
				c.Next()
				return
			}

			// Создаем и устанавливаем куку
			cookie := http.Cookie{
				Name:     "userID",
				Value:    buf.String(),
				Path:     "/",
				MaxAge:   3600 * 24 * 30,
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteLaxMode,
			}

			// Подписываем и устанавливаем куку
			if err := cookies.WriteSigned(c.Writer, cookie, []byte(secretKey)); err != nil {
				log.Printf("Failed to set signed cookie: %v", err)
			} else {
				log.Printf("New cookie set for user: %s", userID)
			}
		}

		c.Next()
	}
}
