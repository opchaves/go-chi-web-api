package app

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/opchaves/go-kom/internal/app/auth/jwt"
	"github.com/opchaves/go-kom/internal/model"
)

type transactionRequest struct {
	Title      string           `json:"title"`
	Note       *string          `json:"note"`
	Amount     float32          `json:"amount"`
	Paid       bool             `json:"paid"`
	TType      string           `json:"t_type"`
	CategoryID uuid.UUID        `json:"category_id"`
	AccountID  uuid.UUID        `json:"account_id"`
	HandledAt  pgtype.Timestamp `json:"handled_at"`
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
		validation.Field(&a.Amount, validation.Required, validation.Min(0)),
		validation.Field(&a.TType, validation.In("expense", "income")),
		validation.Field(&a.CategoryID, validation.Required),
		validation.Field(&a.AccountID, validation.Required),
		validation.Field(&a.HandledAt,
			validation.Required,
			validation.Date(time.DateTime).Min(time.Now()),
		))
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
		ID:     uuid.MustParse(chi.URLParam(r, "id")),
	}

	transaction, err := h.Q.GetTransactionById(r.Context(), params)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	render.Render(w, r, &transactionResponse{transaction})
}

func (h *App) createTransaction(w http.ResponseWriter, r *http.Request) {

}

func (h *App) updateTransaction(w http.ResponseWriter, r *http.Request) {

}

func (h *App) deleteTransaction(w http.ResponseWriter, r *http.Request) {

}
