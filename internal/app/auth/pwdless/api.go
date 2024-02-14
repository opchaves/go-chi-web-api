package pwdless

import (
	"log/slog"
	"net/http"
	"path"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v2"
	"github.com/go-chi/render"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/opchaves/go-chi-web-api/internal/app/auth/jwt"
	"github.com/opchaves/go-chi-web-api/internal/app/email"
	"github.com/opchaves/go-chi-web-api/internal/model"
)

// AuthStorer defines database operations on accounts and tokens.
// type AuthStorer interface {
// 	GetAccount(id string) (*model.User, error)
// 	GetAccountByEmail(email string) (*model.User, error)
// 	UpdateAccount(a *model.User) error
//
// 	GetToken(token string) (*jwt.Token, error)
// 	CreateOrUpdateToken(t *jwt.Token) error
// 	DeleteToken(t *jwt.Token) error
// 	PurgeExpiredToken() error
// }

// Mailer defines methods to send account emails.
type Mailer interface {
	LoginToken(name, email string, c email.ContentLoginToken) error
}

type Resource struct {
	DB        *pgxpool.Pool
	Q         *model.Queries
	LoginAuth *LoginTokenAuth
	TokenAuth *jwt.TokenAuth
	Mailer    Mailer
}

func NewResource(db *pgxpool.Pool, q *model.Queries) (*Resource, error) {
	mailer, err := email.NewMailer()
	if err != nil {
		return nil, err
	}

	loginAuth, err := NewLoginTokenAuth()
	if err != nil {
		return nil, err
	}

	tokenAuth, err := jwt.NewTokenAuth()
	if err != nil {
		return nil, err
	}

	resource := &Resource{
		DB:        db,
		Q:         q,
		LoginAuth: loginAuth,
		TokenAuth: tokenAuth,
		Mailer:    mailer,
	}

	return resource, nil
}

func (rs *Resource) Router() *chi.Mux {
	r := chi.NewRouter()
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Post("/login", rs.login)
	return r
}

type loginRequest struct {
	Email string
}

func (body *loginRequest) Bind(r *http.Request) error {
	body.Email = strings.TrimSpace(body.Email)
	body.Email = strings.ToLower(body.Email)

	return validation.ValidateStruct(body,
		validation.Field(&body.Email, validation.Required, is.Email),
	)
}

func log(r *http.Request) slog.Logger {
	return httplog.LogEntry(r.Context())
}

func (rs *Resource) login(w http.ResponseWriter, r *http.Request) {
	oplog := log(r)

	body := &loginRequest{}
	if err := render.Bind(r, body); err != nil {
		oplog.With("email", body.Email).Warn(err.Error())
		render.Render(w, r, ErrUnauthorized(ErrInvalidLogin))
		return
	}

	user, err := rs.Q.GetUserByEmail(r.Context(), body.Email)
	if err != nil {
		oplog.With("email", body.Email).Warn(err.Error())
		render.Render(w, r, ErrUnauthorized(ErrUnknownLogin))
		return
	}

	// TODO can login

	lt := rs.LoginAuth.CreateToken(user.ID.String())

	go func() {
		content := email.ContentLoginToken{
			Email:  user.Email,
			Name:   user.FirstName,
			URL:    path.Join(rs.LoginAuth.loginURL, lt.Token),
			Token:  lt.Token,
			Expiry: lt.Expiry,
		}
		if err := rs.Mailer.LoginToken(user.FirstName, user.Email, content); err != nil {
			oplog.With("module", "email").Error(err.Error())
		}
	}()

	render.Respond(w, r, http.NoBody)
}
