package app

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/httplog/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/opchaves/go-chi-web-api/internal/model"
)

type App struct {
	DB *pgxpool.Pool
	Q  *model.Queries
}

func New(db *pgxpool.Pool) *App {
	return &App{
		DB: db,
		Q:  model.New(db),
	}
}

func log(r *http.Request) slog.Logger {
	return httplog.LogEntry(r.Context())
}
