package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v2"
	"github.com/opchaves/go-chi-web-api/internal/server"
	"github.com/opchaves/go-chi-web-api/internal/web"
)

func AddRoutes(r *server.Server) error {

	r.Use(middleware.RequestID)
	r.Use(httplog.RequestLogger(r.Logger))
	r.Use(middleware.Heartbeat("/ping"))
	r.Use(middleware.Recoverer)

	r.Route("/", func(r chi.Router) {
		r.NotFound(web.Handler)
	})

	r.Get("/api", func(w http.ResponseWriter, r *http.Request) {
		oplog := httplog.LogEntry(r.Context())
		oplog.Info("api root")
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "Hello, world!"}`))
	})

	r.Get("/warn", func(w http.ResponseWriter, r *http.Request) {
		oplog := httplog.LogEntry(r.Context())
		oplog.Warn("warn here")
		w.WriteHeader(400)
		w.Write([]byte("warn here"))
	})

	return nil
}
