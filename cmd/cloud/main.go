package main

import (
	"fmt"
	"net/http"
	"os"

	conn "github.com/fehmicorp/go/pkg/v1/http"
)

var (
	DefaultPort   = "8080"
	DefaultHost   = "0.0.0.0"
	DefaultAppDir = "app"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = DefaultPort
	}
	host := os.Getenv("HOST")
	if host == "" {
		host = DefaultHost
	}
	dir := os.Getenv("APPDIR")
	if dir == "" {
		dir = DefaultAppDir
	}

	appDir := os.DirFS(dir)
	conn.HandleStaticRouter("/", appDir)

	addr := fmt.Sprintf("%s:%s", host, port)
	fmt.Printf("🚀 Server running on http://%s:%s/\n", host, port)

	if err := http.ListenAndServe(addr, nil); err != nil {
		fmt.Printf("❌ Server failed: %v\n", err)
	}
}
