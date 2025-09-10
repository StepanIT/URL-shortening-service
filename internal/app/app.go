package app

import (
	"database/sql"
	"encoding/gob"
	"fmt"
	"log"

	"github.com/StepanIT/URL-shortening-service/internal/config"
	"github.com/StepanIT/URL-shortening-service/internal/middleware"
	"github.com/StepanIT/URL-shortening-service/internal/server"
	"github.com/StepanIT/URL-shortening-service/internal/storage"
	_ "github.com/lib/pq"
)

func Run() error {
	// load the config
	cfg := config.NewConfig()
	log.Printf("Starting with config: %+v", cfg)

	var (
		db   *sql.DB
		repo storage.URLShortenerRepositories
		err  error
	)

	// подключаем PostgreSQL
	if cfg.DatabaseDsn != "" {
		db, err = sql.Open("postgres", cfg.DatabaseDsn)
		if err != nil {
			log.Printf("Failed to open PostgreSQL: %v", err)
		} else if pingErr := db.Ping(); pingErr != nil {
			log.Printf("Failed to ping PostgreSQL: %v", pingErr)
			db = nil
		} else {
			log.Println("Connected to PostgreSQL successfully")
			repo, err = storage.NewPostgresStorage(cfg.DatabaseDsn)
			if err != nil {
				return fmt.Errorf("failed to init PostgreSQL storage: %w", err)
			}
		}
	}

	// если нет PostgreSQL подключаем FileStorage
	if repo == nil && cfg.FileStoragePath != "" {
		log.Println("Trying File storage...")
		repo, err = storage.NewFileStorage(cfg.FileStoragePath)
		if err != nil {
			return fmt.Errorf("failed to init file storage: %w", err)
		}
		log.Println("Using file storage:", cfg.FileStoragePath)
	}

	// если нет PostgreSQL и FileStorage используем InMemory
	if repo == nil {
		log.Println("Using in-memory storage")
		repo = storage.NewInMemoryStorage()
	}

	u := &middleware.User{}
	gob.Register(u)

	log.Printf("Starting server on %s, BaseURL: %s, Storage: %T", cfg.ServerAddress, cfg.BaseURL, repo)

	// launch the server with all dependencies
	err = server.StartServer(repo, db, cfg.BaseURL, cfg.ServerAddress, cfg.SecretKey)
	if err != nil {
		return fmt.Errorf("server error: %w", err)
	}
	return nil
}
