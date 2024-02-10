package config

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/opchaves/go-chi-web-api/internal/model"
)

const CTX_VERSION = "apiVersion"

type App struct {
	DB *pgxpool.Pool
	Q  *model.Queries
}
