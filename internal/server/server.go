package server

import (
	"github.com/StepanIT/URL-shortening-service/internal/handlers"
	"github.com/StepanIT/URL-shortening-service/internal/storage"
	"github.com/gin-gonic/gin"
)

func StartServer(repo storage.URLShortenerRepositories, baseURL string, serverAddress string) error {

	// pass the selected storage and config to the handler
	h := &handlers.Handler{
		Repo:          repo,
		BaseURL:       baseURL,
		ServerAddress: serverAddress,
	}

	// setting up GIN routes
	router := gin.Default()

	router.POST("/", h.PostHandler)
	router.GET("/get/:id", h.GetHandler)
	router.POST("/api/shorten", h.PostShortenHandler)

	// starting the server
	return router.Run(serverAddress)
}
