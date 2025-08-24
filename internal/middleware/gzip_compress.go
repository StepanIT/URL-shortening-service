package middleware

import (
	"compress/gzip"
	"strings"

	"github.com/gin-gonic/gin"
)

type gzipWriter struct {
	gin.ResponseWriter
	Writer *gzip.Writer
}

func (g *gzipWriter) Write(data []byte) (int, error) {
	return g.Writer.Write(data)
}

func GzipCompress() gin.HandlerFunc {
	return func(c *gin.Context) {
		// проверяем заголовок "Accept-Encoding"
		if !strings.Contains(c.GetHeader("Accept-Encoding"), "gzip") {
			// не поддерживает gzip, пропускаем
			c.Next()
			return
		}

		// устанавливаем заголовок Content-Encoding: gzip
		c.Header("Content-Encoding", "gzip")

		// создаем gzip.Writer для сжатия данных
		gz := gzip.NewWriter(c.Writer)
		defer gz.Close()

		// оборачиваем стандартный ResponseWriter в наш gzipWriter
		gw := &gzipWriter{
			ResponseWriter: c.Writer,
			Writer:         gz,
		}
		c.Writer = gw

		c.Next()
	}
}
