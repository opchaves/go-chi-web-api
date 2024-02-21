package app

import (
	"context"
	"math/big"
	"net/http"

	"github.com/go-chi/render"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/opchaves/go-chi-web-api/internal/app/auth/jwt"
	"github.com/opchaves/go-chi-web-api/internal/model"
	"github.com/opchaves/go-chi-web-api/pkg/pagination"
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

	tx, err := h.DB.Begin(r.Context())
	if err != nil {
		render.Render(w, r, ErrInternal)
	}

	qTx := h.Q.WithTx(tx)
	defer tx.Rollback(r.Context())

	workspace, err := qTx.CreateWorkspace(r.Context(), *input.CreateWorkspaceParams)
	if err != nil {
		oplog.Error("error creating workspace", err)
		render.Render(w, r, ErrRender(err))
		return
	}

	err = createCategories(r.Context(), qTx, workspace)
	if err != nil {
		oplog.Error("error creating categories", err)
		render.Render(w, r, ErrInternal)
	}

	err = createAccounts(r.Context(), qTx, workspace)
	if err != nil {
		oplog.Error("error creating accounts", err)
		render.Render(w, r, ErrInternal)
	}

	_ = tx.Commit(r.Context())

	render.Render(w, r, &createResponse{workspace})
}

func createCategories(ctx context.Context, db *model.Queries, workspace *model.Workspace) error {
	names := []model.CreateCategoryParams{
		{Name: "food", CatType: "expense"},
		{Name: "salary", CatType: "incom"},
		{Name: "health", CatType: "expense"},
		{Name: "transport", CatType: "expense"},
	}

	for _, c := range names {
		c.WorkspaceID = workspace.ID
		c.UserID = workspace.UserID
		if err := db.CreateCategory(ctx, c); err != nil {
			return err
		}
	}

	return nil
}

func createAccounts(ctx context.Context, db *model.Queries, workspace *model.Workspace) error {
	zero := &pgtype.Numeric{Int: big.NewInt(0), Exp: 0, Valid: true}

	names := []model.CreateAccountParams{
		{Name: "checking", Balance: *zero},
		{Name: "savings", Balance: *zero},
	}

	for _, a := range names {
		a.WorkspaceID = workspace.ID
		a.UserID = workspace.UserID
		if _, err := db.CreateAccount(ctx, a); err != nil {
			return err
		}
	}

	return nil
}
