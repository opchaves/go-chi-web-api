package server

import (
	"context"
	"html/template"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httplog/v2"
	"github.com/go-chi/render"
	"github.com/opchaves/go-kom/config"
	"github.com/opchaves/go-kom/server/auth/jwt"
	"github.com/opchaves/go-kom/server/auth/pwdless"
	"github.com/opchaves/go-kom/web"
)

func AddRoutes(r *Server) error {
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.Timeout(15 * time.Second))
	r.Use(httplog.RequestLogger(r.Logger))
	r.Use(render.SetContentType(render.ContentTypeJSON))

	if config.Origins != "" {
		r.Use(corsConfig().Handler)
	}

	api := NewHandler(r.DB, r.Services)

	// TODO: refactor pwdless to use services
	authResource, err := pwdless.NewResource(api.DB, api.Q)

	if err != nil {
		return err
	}

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		t := template.Must(template.ParseFiles("internal/views/index.html"))
		t.Execute(w, nil)
	})

	// Serve react app
	web.FileServer(r, "/app")

	r.Mount("/auth", authResource.Router())
	r.Group(func(r chi.Router) {
		r.Use(authResource.TokenAuth.Verifier())
		r.Use(jwt.Authenticator)
		r.Route("/api/v1", func(r chi.Router) {
			r.Use(apiVersionCtx("v1"))

			r.Get("/workspaces", api.getWorkspace)
			r.Post("/workspaces", api.createWorkspace)
			r.Put("/workspaces/{workspaceID}", api.updateWorkspace)
			r.Delete("/workspaces/{workspaceID}", api.deleteWorkspace)

			r.Get("/transactions", api.getTransactions)
			r.Post("/transactions", api.createTransaction)
			r.Get("/transactions/{transactionID}", api.getTransaction)
		})
	})

	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("ok"))
	})

	return nil
}

func apiVersionCtx(version string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			httplog.LogEntrySetField(ctx, "api_version", slog.StringValue(version))
			r = r.WithContext(context.WithValue(ctx, config.CtxVersion, version))
			next.ServeHTTP(w, r)
		})
	}
}

func corsConfig() *cors.Cors {
	// Basic CORS
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	return cors.New(cors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           86400, // Maximum value not ignored by any of major browsers
	})
}
