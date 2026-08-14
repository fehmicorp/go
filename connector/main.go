package main

import (
	"fmt"
	"log"
	"os"
	"time"

	conn "github.com/fehmicorp/go/pkg/v1/http"
	"github.com/fehmicorp/go/pkg/v1/logger"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Initialize local SQLite logger directory: data/logs.db
	err := logger.Init("./data", "logs.db", 5000, 50, 1*time.Second)
	if err != nil {
		log.Fatalf("Failed to initialize local SQLite logger: %v", err)
	}
	defer logger.Close()

	// Log application boot
	logger.Info("SYSTEM", "Application starting up", map[string]interface{}{
		"version": defConf.Version,
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	host := os.Getenv("HOST")
	if host == "" {
		host = "localhost"
	}
	dir := os.Getenv("APPDIR")
	if dir == "" {
		dir = "app"
	}

	appDir := os.DirFS(dir)
	conn.HandleStaticRouter("/", appDir)
	defConf.Webportal = fmt.Sprintf("http://%s:%s/", host, port)

	initTrayMenu()
	switch Target.OS {
	case "windows":
		Windows()
	case "linux":
		Linux()
	case "darwin":
		Darwin()
	default:
		log.Fatalf("Unsupported operating system: %s", Target.OS)
	}
}
