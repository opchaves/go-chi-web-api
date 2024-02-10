package server

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/go-chi/httplog/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/opchaves/go-chi-web-api/internal/config"
	"github.com/opchaves/go-chi-web-api/internal/model"
)

// Allows to specify options to the server.
type Option func(*Server)

func UsePort(port string) Option {
	return func(s *Server) {
		s.port = port
	}
}

func UseHost(host string) Option {
	return func(s *Server) {
		s.host = host
	}
}

// TODO add other config options like: `version`, ...
func UseLogger(serviceName string) Option {
	logLevel := slog.LevelDebug
	if config.IsProduction {
		logLevel = slog.LevelInfo
	}

	return func(s *Server) {
		s.Logger = httplog.NewLogger(serviceName, httplog.Options{
			JSON:             config.IsProduction,
			LogLevel:         logLevel,
			Concise:          !config.IsProduction,
			RequestHeaders:   config.IsProduction,
			MessageFieldName: "message",
			// Tags: map[string]string{
			// 	"version": "v1.0-81aa4244d9fc8076a",
			// 	"env":     config.Env,
			// },
			QuietDownRoutes: []string{
				"/",
				"/ping",
			},
			QuietDownPeriod: 10 * time.Second,
			// SourceFieldName: "source",
		})
	}
}

func UseDB(ctx context.Context) Option {
	return func(s *Server) {
		db, err := pgxpool.New(ctx, config.DatabaseURL)
		if err != nil {
			s.Logger.ErrorContext(ctx, "Could not connect to database", err)
			os.Exit(1)
		}

		s.Logger.DebugContext(ctx, "Database successfully connected.")

		s.DB = db
		s.Q = model.New(db)
	}
}
