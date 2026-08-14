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
	err := logger.Init("./data", "logs.db", 5000, 50, 1*time.Second)
	if err != nil {
		log.Fatalf("Failed to initialize local SQLite logger: %v", err)
	}
	defer logger.Close()
	Cfg := LoadConfig()
	logger.Info("SYSTEM", "Application starting up", map[string]interface{}{
		"version":    defConf.Version,
		"at_startup": Cfg.AtStartup,
	})
	appDir := os.DirFS(Cfg.AppDir)
	fmt.Println(appDir)
	conn.HandleStaticRouter("/", appDir)
	defConf.Webportal = fmt.Sprintf("http://%s:%s", Cfg.Host, Cfg.Port)
	serverMutex.Lock()
	if Cfg.AtStartup {
		fmt.Println("Auto-starting web server due to configuration (at_startup = true)...")
		StartServerLocked()
	}
	serverMutex.Unlock()
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
