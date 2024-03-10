package pwdless

import (
	"fmt"
	"log/slog"
	"net/http"
	"path"
	"strings"

	"github.com/go-chi/httplog/v2"
	"github.com/go-chi/render"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/opchaves/go-kom/server/email"
)

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
	fmt.Println("login", body.Email)

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
