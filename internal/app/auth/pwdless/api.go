package pwdless

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/opchaves/go-kom/internal/app/auth/jwt"
	"github.com/opchaves/go-kom/internal/app/email"
	"github.com/opchaves/go-kom/internal/model"
)

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
	r.Post("/token", rs.token)
	r.Group(func(r chi.Router) {
		r.Use(rs.TokenAuth.Verifier())
		r.Use(jwt.AuthenticateRefreshJWT)
		r.Post("/refresh", rs.refresh)
		r.Post("/logout", rs.logout)
	})
	return r
}
