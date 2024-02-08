package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/opchaves/go-chi-web-api/internal/server"
	"github.com/opchaves/go-chi-web-api/internal/web"
)

func AddRoutes(r *server.Server) error {
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/", func(r chi.Router) {
		r.NotFound(web.Handler)
	})

	r.Route("/api", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"message": "Hello, world!"}`))
		})
	})

	return nil
}
