package services

import (
	"github.com/opchaves/go-kom/internal/stores"
)

type Services struct {
	Workspace WorkspaceService
}

func New(st *stores.Stores) *Services {
	return &Services{
		Workspace: &workspaceService{stores: st},
	}
}
