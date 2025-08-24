package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func GzipDecompress() gin.HandlerFunc {
	return func(c *gin.Context) {
		// проверяем заголовок Content-Encoding
		if strings.Contains(c.GetHeader("Content-Encoding"), "gzip") {
			// создаем gzip.NewReader для чтения сжатого тела запроса
			gz, err := gzip.NewReader(c.Request.Body)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"error": "failed to read gzip body",
				})
				return
			}
			defer gz.Close()

			// читаем разжатые данные
			body, err := io.ReadAll(gz)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"error": "failed to decompress body",
				})
				return
			}

			// заменяем тело запроса на разжатое
			c.Request.Body = io.NopCloser(strings.NewReader(string(body)))

			// удаляем заголовок Content-Encoding
			c.Request.Header.Del("Content-Encoding")
		}
		c.Next()
	}
}
