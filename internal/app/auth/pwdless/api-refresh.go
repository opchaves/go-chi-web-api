package pwdless

import (
	"net/http"
	"time"

	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/opchaves/go-kom/internal/app/auth/jwt"
	"github.com/opchaves/go-kom/config"
	"github.com/opchaves/go-kom/model"
)

func (rs *Resource) refresh(w http.ResponseWriter, r *http.Request) {
	oplog := log(r)
	rt := jwt.RefreshTokenFromCtx(r.Context())

	token, err := rs.Q.GetToken(r.Context(), rt)
	if err != nil {
		render.Render(w, r, ErrUnauthorized(jwt.ErrTokenExpired))
		return
	}

	// if token is expired, delete it and return unauthorized
	if time.Now().After(token.ExpiresAt.Time) {
		_ = rs.Q.DeleteTokenByID(r.Context(), token.ID)
		render.Render(w, r, ErrUnauthorized(jwt.ErrTokenExpired))
		return
	}

	user, err := rs.Q.GetUserById(r.Context(), token.UserID)
	if err != nil {
		render.Render(w, r, ErrUnauthorized(ErrUnknownLogin))
		return
	}

	if !user.CanLogin() {
		render.Render(w, r, ErrUnauthorized(ErrLoginDisabled))
		return
	}

	var expiresAt pgtype.Timestamp
	expiresAt.Scan(time.Now().Add(config.JwtRefreshExpiry))
	tokenParams := model.UpdateTokenParams{
		ID:        token.ID,
		ExpiresAt: expiresAt,
		Token:     uuid.Must(uuid.NewRandom()).String(),
	}

	claims := jwt.AppClaims{
		ID:    user.ID.String(),
		Sub:   user.FirstName,
		Roles: []string{"user"},
	}

	refreshClaims := jwt.RefreshClaims{
		ID:    int(token.ID),
		Token: tokenParams.Token,
	}

	access, refresh, err := rs.TokenAuth.GenTokenPair(claims, refreshClaims)
	if err != nil {
		oplog.Error(err.Error())
		render.Render(w, r, ErrInternalServerError)
		return
	}

	if err := rs.Q.UpdateToken(r.Context(), tokenParams); err != nil {
		oplog.Error(err.Error())
		render.Render(w, r, ErrInternalServerError)
		return
	}

	render.Respond(w, r, &tokenResponse{
		Access:  access,
		Refresh: refresh,
	})
	// TODO: add `last_login` to user???
	// TODO use int `id` for tokens instead
}
