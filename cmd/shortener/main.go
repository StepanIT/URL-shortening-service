package main

import (
	"flag"

	"github.com/StepanIT/URL-shortening-service/cmd/shortener/server"
)

func main() {
	addr := flag.String("a", "", "server address")
	base := flag.String("b", "", "base URL")
	file := flag.String("f", "", "file storage path")
	flag.Parse()

	server.Handler(*addr, *base, *file)
}
