package app

import (
	"net/http"

	"github.com/go-chi/render"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/opchaves/go-kom/internal/app/auth/jwt"
	"github.com/opchaves/go-kom/internal/model"
	"github.com/opchaves/go-kom/pkg/pagination"
)

type createResponse struct {
	*model.Workspace
}

func (rd *createResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type createInput struct {
	*model.CreateWorkspaceParams
}

func (a *createInput) Bind(r *http.Request) error {
	return validation.ValidateStruct(a,
		validation.Field(&a.Name, validation.Required),
		validation.Field(&a.Currency, validation.In("BRL", "USD", "EUR")),
		validation.Field(&a.Language, validation.In("pt-br", "en-us")),
	)
}

func (h *App) getWorkspace(w http.ResponseWriter, r *http.Request) {
	oplog := log(r)
	oplog.Info("list workspaces by user")

	claims := jwt.ClaimsFromCtx(r.Context())
	pages := pagination.NewFromRequest(r, -1)
	params := model.GetWorkspacesByUserParams{
		UserID: uuid.MustParse(claims.ID),
		Limit:  int32(pages.Limit()),
		Offset: int32(pages.Offset()),
	}

	workspaces, err := h.Q.GetWorkspacesByUser(r.Context(), params)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
	}

	list := []render.Renderer{}
	for _, w := range workspaces {
		list = append(list, &createResponse{w})
	}

	render.RenderList(w, r, list)
}

func (h *App) createWorkspace(w http.ResponseWriter, r *http.Request) {
	oplog := log(r)
	oplog.Info("create new workspace")

	input := &createInput{}
	if err := render.Bind(r, input); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	claims := jwt.ClaimsFromCtx(r.Context())
	input.UserID = uuid.MustParse(claims.ID)

	workspace, err := h.Services.Workspace.Create(r.Context(), input.CreateWorkspaceParams)

	if err != nil {
		oplog.Error("error creating workspace", err)
		render.Render(w, r, ErrRender(err))
		return
	}
	render.Render(w, r, &createResponse{workspace})
}
