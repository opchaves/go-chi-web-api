package services

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/opchaves/go-kom/model"
)

type Services struct {
	DB *pgxpool.Pool
	Q  *model.Queries

	Workspace WorkspaceService
}

func New(db *pgxpool.Pool) *Services {
	q := model.New(db)
	return &Services{
		DB: db,
		Q:  q,

		Workspace: &workspaceService{DB: db, Q: q},
	}
}
