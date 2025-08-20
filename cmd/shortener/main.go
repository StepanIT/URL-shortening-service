package main

import (
	"log"

	"github.com/StepanIT/URL-shortening-service/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatal("Application failed to start: ", err)
	}
}
