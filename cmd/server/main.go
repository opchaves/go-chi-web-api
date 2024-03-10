package main

import (
	"context"
	"os"

	"github.com/opchaves/go-kom/config"
	"github.com/opchaves/go-kom/internal/app"
	"github.com/opchaves/go-kom/server"
)

func main() {
	ctx := context.TODO()

	s := server.New(
		config.Name,
		server.UseHost(config.Host),
		server.UsePort(config.Port),
		server.UseLogger(),
		server.UseDB(ctx),
		server.UseServices(),
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
