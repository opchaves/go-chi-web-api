package pwdless

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/opchaves/go-kom/internal/app/auth/jwt"
)

func (rs *Resource) logout(w http.ResponseWriter, r *http.Request) {
	rt := jwt.RefreshTokenFromCtx(r.Context())
	token, err := rs.Q.GetToken(r.Context(), rt)
	if err != nil {
		render.Render(w, r, ErrUnauthorized(jwt.ErrTokenExpired))
		return
	}
	rs.Q.DeleteTokenByID(r.Context(), token.ID)

	render.Respond(w, r, http.NoBody)
}
