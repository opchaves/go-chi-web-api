package jwt

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/go-chi/httplog/v2"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/opchaves/go-kom/internal/config"
)

// ClaimsFromCtx retrieves the parsed AppClaims from request context.
func ClaimsFromCtx(ctx context.Context) AppClaims {
	return ctx.Value(config.CtxClaims).(AppClaims)
}

// UserIDFromCtx retrieves the user ID from the parsed AppClaims in request context.
func UserIDFromCtx(ctx context.Context) uuid.UUID {
	claims := ClaimsFromCtx(ctx)
	return uuid.MustParse(claims.ID)
}

// RefreshTokenFromCtx retrieves the parsed refresh token from context.
func RefreshTokenFromCtx(ctx context.Context) string {
	return ctx.Value(config.CtxRefreshToken).(string)
}

func log(r *http.Request) slog.Logger {
	return httplog.LogEntry(r.Context())
}

// Authenticator is a default authentication middleware to enforce access from the
// Verifier middleware request context values. The Authenticator sends a 401 Unauthorized
// response for any unverified tokens and passes the good ones through.
func Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		oplog := log(r)

		token, claims, err := jwtauth.FromContext(r.Context())
		if err != nil {
			oplog.Warn(err.Error())
			render.Render(w, r, ErrUnauthorized(ErrTokenUnauthorized))
			return
		}

		if err := jwt.Validate(token); err != nil {
			render.Render(w, r, ErrUnauthorized(ErrTokenExpired))
			return
		}

		// Token is authenticated, parse claims
		var c AppClaims
		err = c.ParseClaims(claims)
		if err != nil {
			oplog.Error(err.Error())
			render.Render(w, r, ErrUnauthorized(ErrInvalidAccessToken))
			return
		}

		// Set AppClaims on context
		ctx := context.WithValue(r.Context(), config.CtxClaims, c)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AuthenticateRefreshJWT checks validity of refresh tokens and is only used for access token refresh and logout requests. It responds with 401 Unauthorized for invalid or expired refresh tokens.
func AuthenticateRefreshJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		oplog := log(r)

		token, claims, err := jwtauth.FromContext(r.Context())
		if err != nil {
			oplog.Warn(err.Error())
			render.Render(w, r, ErrUnauthorized(ErrTokenUnauthorized))
			return
		}

		if err := jwt.Validate(token); err != nil {
			render.Render(w, r, ErrUnauthorized(ErrTokenExpired))
			return
		}

		// Token is authenticated, parse refresh token string
		var c RefreshClaims
		err = c.ParseClaims(claims)
		if err != nil {
			oplog.Error(err.Error())
			render.Render(w, r, ErrUnauthorized(ErrInvalidRefreshToken))
			return
		}

		// Set refresh token string on context
		ctx := context.WithValue(r.Context(), config.CtxRefreshToken, c.Token)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
