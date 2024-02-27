package web

import (
	"embed"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"path"
	"strings"

	"github.com/opchaves/go-kom/internal/config"
)

var (
	//go:embed all:build
	assets embed.FS

	assetsDir = "build"
	errDir    = errors.New("path is dir")
	maxAge    = 604800 // 7 days
)

func WebHandler(w http.ResponseWriter, r *http.Request) {
	if strings.HasSuffix(r.URL.Path, "index.html") {
		newPath := strings.TrimSuffix(r.URL.Path, "index.html")
		http.Redirect(w, r, newPath, http.StatusMovedPermanently)
		return
	}

	// try serving an assets
	err := readFile(r.URL.Path, w)
	if err == nil {
		return
	}

	reqPath := ""
	if err != nil {
		if err != errDir {
			_ = readFile("index.html", w)
			return
		}
		if r.URL.Path != "/" {
			reqPath = r.URL.Path
		}
	}
	reqPath = path.Join(reqPath, "index.html")

	_ = readFile(reqPath, w)
}

func readFile(reqPath string, w http.ResponseWriter) error {
	pathname := path.Join(assetsDir, reqPath)
	f, err := assets.Open(pathname)
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
