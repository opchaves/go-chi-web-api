package app

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/httplog/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/opchaves/go-kom/model"
	"github.com/opchaves/go-kom/services"
	"github.com/opchaves/go-kom/stores"
)

type App struct {
	DB       *pgxpool.Pool
	Q        *model.Queries
	Stores   *stores.Stores
	Services *services.Services
}

func New(db *pgxpool.Pool, stores *stores.Stores, servs *services.Services) *App {
	return &App{
		DB:       db,
		Q:        model.New(db),
		Stores:   stores,
		Services: servs,
	}
}

func log(r *http.Request) slog.Logger {
	return httplog.LogEntry(r.Context())
}
