package web

import (
	"embed"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/opchaves/go-kom/config"
)

// TODO use file server implementation from chi examples

var (
	//go:embed all:dist/*
	assetsFS embed.FS
	// errDir    = errors.New("path is dir")
	// maxAge    = 604800 // 7 days
)

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func FileServer(r chi.Router, path string) {
	root := getFileSystem()

	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}

func getFileSystem() http.FileSystem {
	if config.IsProduction {
		staticContent, err := fs.Sub(fs.FS(assetsFS), "dist")
		if err != nil {
			panic(err)
		}
		return http.FS(staticContent)
	}

	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "web", "dist"))
	return filesDir
}
