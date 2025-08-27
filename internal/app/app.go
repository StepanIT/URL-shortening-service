package app

import (
	"encoding/gob"
	"fmt"
	"log"

	"github.com/StepanIT/URL-shortening-service/internal/config"
	"github.com/StepanIT/URL-shortening-service/internal/handlers"
	"github.com/StepanIT/URL-shortening-service/internal/server"
	"github.com/StepanIT/URL-shortening-service/internal/storage"
)

func Run() error {
	// load the config
	cfg := config.NewConfig()
	log.Printf("Starting with config: %+v", cfg)

	// path to file storage
	var repo storage.URLShortenerRepositories
	var err error
	if cfg.FileStoragePath != "" {
		// use FileStorage
		repo, err = storage.NewFileStorage(cfg.FileStoragePath)
		if err != nil {
			return fmt.Errorf("config error: %w", err)
		}
		log.Println("Using file storage:", cfg.FileStoragePath)
	} else {
		// use in-memory storage
		repo = storage.NewInMemoryStorage()
		log.Println("Using in-memory storage")
	}
	name := "Ignat"
	u := &handlers.User{
		Name: name,
	}

	gob.Register(u)
	log.Printf("Starting server on %s, %s, %s, %s", cfg.ServerAddress, cfg.BaseURL, repo, cfg.SecretKey)

	// launch the server with all dependencies
	err = server.StartServer(repo, cfg.BaseURL, cfg.ServerAddress, cfg.SecretKey)
	if err != nil {
		return fmt.Errorf("server error: %w", err)
	}
	return nil
}
