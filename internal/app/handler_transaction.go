package app

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/google/uuid"
	"github.com/opchaves/go-kom/internal/app/auth/jwt"
	"github.com/opchaves/go-kom/internal/model"
	"github.com/opchaves/go-kom/pkg/util"
)

type transactionRequest struct {
	Title       string    `json:"title"`
	Note        *string   `json:"note"`
	Amount      string    `json:"amount"`
	Paid        bool      `json:"paid"`
	TType       string    `json:"t_type"`
	HandledAt   string    `json:"handled_at"`
	WorkspaceID uuid.UUID `json:"workspace_id"`
	UserID      uuid.UUID `json:"user_id"`
	CategoryID  uuid.UUID `json:"category_id"`
	AccountID   uuid.UUID `json:"account_id"`
}

type transactionResponse struct {
	*model.Transaction
}

func (rd *transactionResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (a *transactionRequest) Bind(r *http.Request) error {
	return validation.ValidateStruct(a,
		validation.Field(&a.Title, validation.Required),
		validation.Field(&a.Amount, validation.Required, is.Float),
		validation.Field(&a.TType, validation.In("expense", "income")),
		validation.Field(&a.CategoryID, validation.Required),
		validation.Field(&a.AccountID, validation.Required),
		validation.Field(&a.HandledAt, validation.Required, validation.Date(time.RFC3339)),
		validation.Field(&a.CategoryID, validation.Required, is.UUID),
		validation.Field(&a.AccountID, validation.Required, is.UUID),
	)
}

func (h *App) getTransactions(w http.ResponseWriter, r *http.Request) {
	oplog := log(r)
	oplog.Info("list transactions by user")

	userId := jwt.UserIDFromCtx(r.Context())

	transactions, err := h.Q.GetTransactionsByUser(r.Context(), userId)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	list := []render.Renderer{}
	for _, w := range transactions {
		list = append(list, &transactionResponse{w})
	}

	render.RenderList(w, r, list)
}

func (h *App) getTransaction(w http.ResponseWriter, r *http.Request) {
	oplog := log(r)
	oplog.Info("get transaction by user")

	params := model.GetTransactionByIdParams{
		UserID: jwt.UserIDFromCtx(r.Context()),
		ID:     uuid.MustParse(chi.URLParam(r, "transactionID")),
	}

	transaction, err := h.Q.GetTransactionById(r.Context(), params)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	render.Render(w, r, &transactionResponse{transaction})
}

func (h *App) createTransaction(w http.ResponseWriter, r *http.Request) {
	oplog := log(r)
	oplog.Info("creating new transaction")

	input := &transactionRequest{}
	if err := render.Bind(r, input); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	amount := util.ParseNumeric(input.Amount)
	handledAt := util.ParseTimestamp(input.HandledAt)
	input.UserID = jwt.UserIDFromCtx(r.Context())
	input.WorkspaceID = jwt.WorkspaceIDFromCtx(r.Context())

	if amount == nil || handledAt == nil {
		render.Render(w, r, ErrInvalidRequest(fmt.Errorf("failed to parse input")))
		return
	}

	txParams := model.CreateTransactionParams{
		Title:       input.Title,
		Note:        input.Note,
		Amount:      *amount,
		Paid:        input.Paid,
		TType:       input.TType,
		HandledAt:   *handledAt,
		WorkspaceID: input.WorkspaceID,
		UserID:      input.UserID,
		CategoryID:  input.CategoryID,
		AccountID:   input.AccountID,
	}

	tx, err := h.Q.CreateTransaction(r.Context(), txParams)
	if err != nil {
		oplog.Error("error creating transaction", err)
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	render.Render(w, r, &transactionResponse{tx})
}

func (h *App) updateTransaction(w http.ResponseWriter, r *http.Request) {

}

func (h *App) deleteTransaction(w http.ResponseWriter, r *http.Request) {

}
