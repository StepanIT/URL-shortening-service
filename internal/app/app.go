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

	var db *sql.DB
	var err error
	if cfg.DatabaseDsn != "" {
		db, err := sql.Open("postgres", cfg.DatabaseDsn)
		if err != nil {
			panic(err)
		}
		defer db.Close()

		if err := db.Ping(); err != nil {
			log.Println("Database connection failed:", err)
		} else {
			log.Println("Connected to PostgreSQL successfully")
		}
	}

	// path to file storage
	var repo storage.URLShortenerRepositories

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

	u := &middleware.User{}
	gob.Register(u)

	log.Printf("Starting server on %s, %v, %s, %s, %s", cfg.ServerAddress, db, cfg.BaseURL, repo, cfg.SecretKey)

	// launch the server with all dependencies
	err = server.StartServer(repo, db, cfg.BaseURL, cfg.ServerAddress, cfg.SecretKey)
	if err != nil {
		return fmt.Errorf("server error: %w", err)
	}
	return nil
}
