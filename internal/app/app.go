package app

import (
	"log"

	"github.com/StepanIT/URL-shortening-service/internal/config"
	"github.com/StepanIT/URL-shortening-service/internal/server"
)

func Run() {
	cfg := config.NewConfig()
	log.Printf("Starting with config: %+v", cfg)
	server.Handler(cfg)
}
