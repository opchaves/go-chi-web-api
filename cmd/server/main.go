package main

import (
	"log"
	"os"

	"github.com/opchaves/go-chi-web-api/internal/app"
	"github.com/opchaves/go-chi-web-api/internal/config"
	"github.com/opchaves/go-chi-web-api/internal/server"
)

func main() {
	s := server.New(
		config.Name,
		server.UseHost(config.Host),
		server.UsePort(config.Port),
	)

	if err := app.AddRoutes(s); err != nil {
		log.Fatalf("error adding routes: %v", err)
		os.Exit(1)
	}

	if err := s.Run(); err != nil {
		log.Fatalf("error running server: %v", err)
		os.Exit(1)
	}
}
