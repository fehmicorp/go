package main

import (
	"fmt"
	"log"

	"github.com/fehmicorp/go/pkg/v1/win/systray"
	"github.com/wailsapp/wails/v3/pkg/application"
)

func Darwin() {
	app := application.New(application.Options{
		Name:        defConf.Title,
		Description: defConf.Desc,
	})
	trayIcon, err := AssetsFS.ReadFile("assets/icon.png")
	if err != nil {
		log.Fatalf("failed to read embedded asset `assets/icon.png`: %v", err)
	}
	Tooltip := fmt.Sprintf("%s\nVersion: %s\nStatus: %s", defConf.Title, defConf.Version, serverStatus)
	customMenus := PrepareMenuItems()

	// Initialize tray instance for macOS menu bar
	trayInstance = systray.NewTrayManager(app, trayIcon, Tooltip, customMenus)

	if err := app.Run(); err != nil {
		log.Fatalf("❌ macOS system tray service failed: %v", err)
	}
}
