package main

import (
	"log"

	config "github.com/StepanIT/URL-shortening-service"
	"github.com/StepanIT/URL-shortening-service/cmd/shortener/server"
)

func main() {
	cfg := config.NewConfig()
	log.Printf("Starting with config: %+v", cfg)
	server.Handler(cfg)
}
