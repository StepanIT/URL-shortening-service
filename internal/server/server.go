package server

import (
	"github.com/StepanIT/URL-shortening-service/internal/handlers"
	"github.com/StepanIT/URL-shortening-service/internal/middleware"
	"github.com/StepanIT/URL-shortening-service/internal/storage"
	"github.com/gin-gonic/gin"
)

func StartServer(repo storage.URLShortenerRepositories, baseURL string, serverAddress string, secretKey string) error {

	// pass the selected storage and config to the handler
	h := &handlers.Handler{
		Repo:          repo,
		BaseURL:       baseURL,
		ServerAddress: serverAddress,
		SecretKey:     secretKey,
	}

	// setting up GIN routes
	router := gin.Default()

	router.GET("/set", h.SetCookieHandler)
	router.GET("/get", h.GetCookieHandler)

	router.Use(middleware.GzipDecompress())
	router.Use(middleware.GzipCompress())

	router.POST("/", h.PostHandler)
	router.GET("/get/:id", h.GetHandler)
	router.POST("/api/shorten", h.PostShortenHandler)

	// starting the server
	return router.Run(serverAddress)
}
