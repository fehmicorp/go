package http

import (
	"io/fs"
	"net/http"
	"strings"
)

func HandleStaticRouter(urlPattern string, staticFS fs.FS) {
	fileServer := http.FileServer(http.FS(staticFS))

	http.HandleFunc(urlPattern, func(w http.ResponseWriter, r *http.Request) {
		cleanPath := strings.Trim(r.URL.Path, "/")
		if cleanPath == "" {
			serveFile(w, staticFS, "index.html")
			return
		}
		if cleanPath == "onboard" {
			serveFile(w, staticFS, "onboard.html")
			return
		}
		if _, err := fs.Stat(staticFS, cleanPath); err == nil {
			r.URL.Path = cleanPath
			fileServer.ServeHTTP(w, r)
			return
		}
		htmlPath := cleanPath + ".html"
		if _, err := fs.Stat(staticFS, htmlPath); err == nil {
			r.URL.Path = htmlPath
			fileServer.ServeHTTP(w, r)
			return
		}
		serveFile(w, staticFS, "index.html")
	})
}

func serveFile(w http.ResponseWriter, staticFS fs.FS, filename string) {
	page, err := fs.ReadFile(staticFS, filename)
	if err != nil {
		http.Error(w, filename+" missing in target directory.", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(page)
}
