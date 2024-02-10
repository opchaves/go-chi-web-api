package workspace

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v2"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/opchaves/go-chi-web-api/internal/config"
	"github.com/opchaves/go-chi-web-api/internal/model"
)

type WorkspaceResource struct {
	*config.App
}

func NewWorkspaceResource(app *config.App) *WorkspaceResource {
	return &WorkspaceResource{app}
}

func (wr *WorkspaceResource) Router() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/{id}", wr.get)
	r.Post("/", wr.createWorkspace)

	return r
}

func log(r *http.Request) slog.Logger {
	return httplog.LogEntry(r.Context())
}

// TODO maybe have each route handler as a file

func (wr *WorkspaceResource) get(w http.ResponseWriter, r *http.Request) {
	oplog := log(r)
	oplog.Info("list workspaces by user")

	// TODO get userId from context after jwt auth is implemented
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
	}

	workspaces, err := wr.Q.GetWorkspacesByUser(r.Context(), id)

	list := []render.Renderer{}
	for _, w := range workspaces {
		list = append(list, &createResponse{w})
	}

	render.RenderList(w, r, list)
}

type createInput struct {
	*model.CreateWorkspaceParams
}

func (a *createInput) Bind(r *http.Request) error {
	return nil
}

type createResponse struct {
	*model.Workspace
}

func (rd *createResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (wr *WorkspaceResource) createWorkspace(w http.ResponseWriter, r *http.Request) {
	oplog := log(r)
	oplog.Info("create new workspace")

	input := &createInput{}
	if err := render.Bind(r, input); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	// input.UserID, _ = uuid.NewV7()
	input.UserID, _ = uuid.Parse("48edef85-2738-4295-8555-1930f0c844e9")

	workspace, err := wr.Q.CreateWorkspace(r.Context(), *input.CreateWorkspaceParams)
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	render.Render(w, r, &createResponse{workspace})
}
