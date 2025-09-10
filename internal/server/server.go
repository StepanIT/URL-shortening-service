package server

import (
	"database/sql"

	"github.com/StepanIT/URL-shortening-service/internal/handlers"
	"github.com/StepanIT/URL-shortening-service/internal/middleware"
	"github.com/StepanIT/URL-shortening-service/internal/storage"
	"github.com/gin-gonic/gin"
)

func StartServer(repo storage.URLShortenerRepositories, db *sql.DB, baseURL string, serverAddress string, secretKey string) error {

	// pass the selected storage and config to the handler
	h := &handlers.Handler{
		Repo:          repo,
		BaseURL:       baseURL,
		ServerAddress: serverAddress,
		SecretKey:     secretKey,
	}
	pingHandler := handlers.NewPingHandler(db)
	// setting up GIN routes
	router := gin.Default()

	router.Use(middleware.GzipDecompress())
	router.Use(middleware.GzipCompress())
	router.Use(middleware.Auth(secretKey))

	router.GET("/get/:id", h.GetHandler)
	router.GET("/:id", h.GetHandler)
	router.GET("/api/user/urls", h.GetUserURLsHandler)

	router.POST("/", h.PostHandler)
	router.POST("/api/shorten", h.PostShortenHandler)
	router.POST("/api/shorten/batch", h.PostShortenBatchHandler)

	router.GET("/ping", pingHandler.Ping)

	// starting the server
	return router.Run(serverAddress)
}
