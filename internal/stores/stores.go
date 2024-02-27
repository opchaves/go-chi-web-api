package stores

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/opchaves/go-chi-web-api/internal/model"
)

type Stores struct {
	DB *pgxpool.Pool
	Q  *model.Queries
}

func New(db *pgxpool.Pool) *Stores {
	return &Stores{
		DB: db,
		Q:  model.New(db),
	}
}
