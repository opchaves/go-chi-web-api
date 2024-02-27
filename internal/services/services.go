package services

import (
	"github.com/opchaves/go-chi-web-api/internal/stores"
)

type Services struct {
	Workspace WorkspaceService
}

func New(st *stores.Stores) *Services {
	return &Services{
		Workspace: &workspaceService{stores: st},
	}
}
