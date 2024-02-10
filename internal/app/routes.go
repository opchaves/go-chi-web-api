package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v2"
	"github.com/go-chi/render"
	"github.com/opchaves/go-chi-web-api/internal/app/workspace"
	"github.com/opchaves/go-chi-web-api/internal/config"
	"github.com/opchaves/go-chi-web-api/internal/server"
	"github.com/opchaves/go-chi-web-api/internal/web"
)

func AddRoutes(r *server.Server) error {
	r.Use(middleware.RequestID)
	r.Use(httplog.RequestLogger(r.Logger))
	r.Use(middleware.Heartbeat("/ping"))
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	app := &config.App{
		DB: r.DB,
		Q:  r.Q,
	}

	workspaceResource := workspace.NewWorkspaceResource(app)

	r.Route("/api", func(r chi.Router) {
		r.Get("/", apiHello)
		r.Route("/v1", func(r chi.Router) {
			r.Use(apiVersionCtx("v1"))
			r.Get("/", apiHello)
			r.Mount("/workspaces", workspaceResource.Router())
		})
	})

	r.Get("/*", web.WebHandler)

	return nil
}

func apiVersionCtx(version string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			httplog.LogEntrySetField(ctx, config.CTX_VERSION, slog.StringValue(version))
			r = r.WithContext(context.WithValue(ctx, config.CTX_VERSION, version))
			next.ServeHTTP(w, r)
		})
	}
}

func apiHello(w http.ResponseWriter, r *http.Request) {
	apiVersion := r.Context().Value(config.CTX_VERSION)
	msg := fmt.Sprintf(`Hello from %v`, r.URL.Path)
	result := map[string]interface{}{"message": msg}
	if apiVersion != nil {
		result["apiVersion"] = apiVersion
	}
	render.JSON(w, r, result)
}
