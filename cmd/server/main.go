package main

import (
	"context"
	"os"

	"github.com/opchaves/go-chi-web-api/internal/app"
	"github.com/opchaves/go-chi-web-api/internal/config"
	"github.com/opchaves/go-chi-web-api/internal/server"
)

func main() {
	ctx := context.TODO()

	s := server.New(
		config.Name,
		server.UseLogger(config.Name),
		server.UseDB(ctx),
		server.UseHost(config.Host),
		server.UsePort(config.Port),
	)

	defer s.DB.Close()

	if err := app.AddRoutes(s); err != nil {
		s.Logger.ErrorContext(ctx, "error adding routes", err)
		os.Exit(1)
	}

	if err := s.Run(); err != nil {
		s.Logger.ErrorContext(ctx, "error running server", err)
		os.Exit(1)
	}
}
