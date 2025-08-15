package main

import (
	config "github.com/StepanIT/URL-shortening-service"
	"github.com/StepanIT/URL-shortening-service/cmd/shortener/server"
)

func main() {
	cfg := config.NewConfig()
	server.Handler(cfg)
}
