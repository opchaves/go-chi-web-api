package workspace

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v2"
	"github.com/opchaves/go-chi-web-api/internal/config"
)

type WorkspaceResource struct {
	*config.App
}

func NewWorkspaceResource(app *config.App) *WorkspaceResource {
	return &WorkspaceResource{app}
}

func (wr *WorkspaceResource) Router() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/", wr.get)

	return r
}

func log(r *http.Request) slog.Logger {
	return httplog.LogEntry(r.Context())
}

func (wr *WorkspaceResource) get(w http.ResponseWriter, r *http.Request) {
	oplog := log(r)
	oplog.Info("api root")
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "Hello, world!"}`))
}
