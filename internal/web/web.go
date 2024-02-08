package web

import (
	"embed"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"path"

	"github.com/opchaves/go-chi-web-api/internal/config"
)

var (
	//go:embed all:public
	files embed.FS

	rootPath = "/"
	errDir   = errors.New("path is dir")
	maxAge   = 604800 // 7 days
)

// TODO https://github.com/go-chi/chi/blob/master/_examples/fileserver/main.go

func Handler(w http.ResponseWriter, r *http.Request) {
	err := readFile(files, "public", r.URL.Path, w)
	if err == nil {
		return
	}

	reqDir := "public"
	if err != nil {
		if err != errDir {
			// TODO render 404 page instead???
			http.Error(w, "404 page not found", http.StatusNotFound)
			return
		}
		if r.URL.Path != rootPath {
			reqDir = path.Join("public", r.URL.Path)
		}
	}

	err = readFile(files, reqDir, "index.html", w)
	if err != nil {
		// TODO render 404 page instead???
		http.Error(w, "404 page not found", http.StatusNotFound)
		return
	}
}

func readFile(fs embed.FS, prefix, reqPath string, w http.ResponseWriter) error {
	pathname := path.Join(prefix, reqPath)
	f, err := fs.Open(pathname)
	if err != nil {
		return err
	}
	defer f.Close()

	stat, _ := f.Stat()
	if stat.IsDir() {
		return errDir
	}

	contentType := mime.TypeByExtension(path.Ext(reqPath))
	w.Header().Set("Content-Type", contentType)

	if config.IsProduction {
		cacheControl := fmt.Sprintf("public, max-age=%d, immutable", maxAge)
		w.Header().Set("Cache-Control", cacheControl)
	}

	_, err = io.Copy(w, f)
	return err
}
