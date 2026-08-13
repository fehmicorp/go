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

		// 1. Explicitly serve index.html for root path
		if cleanPath == "" {
			serveFile(w, staticFS, "index.html")
			return
		}

		// 2. Explicitly serve onboard.html for onboard path
		if cleanPath == "onboard" {
			serveFile(w, staticFS, "onboard.html")
			return
		}

		// 3. Check if physical asset exists (e.g., _next/static/..., favicon.ico, images/logo)
		if _, err := fs.Stat(staticFS, cleanPath); err == nil {
			r.URL.Path = cleanPath
			fileServer.ServeHTTP(w, r)
			return
		}

		// 4. Fallback for generic file mapping or dynamic sub-routes
		htmlPath := cleanPath + ".html"
		if _, err := fs.Stat(staticFS, htmlPath); err == nil {
			r.URL.Path = htmlPath
			fileServer.ServeHTTP(w, r)
			return
		}

		// 5. Final fallback to root index.html
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
