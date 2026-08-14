package main

import (
	"fmt"
	"log"

	"github.com/fehmicorp/go/pkg/v1/win/systray"
	"github.com/wailsapp/wails/v3/pkg/application"
)

func Linux() {
	app := application.New(application.Options{
		Name:        defConf.Title,
		Description: defConf.Desc,
	})
	trayIcon, err := AssetsFS.ReadFile("assets/icon.png")
	if err != nil {
		log.Fatalf("failed to read embedded asset `assets/icon.png`: %s", err)
	}
	Tooltip := fmt.Sprintf("%s\nVersion: %s\nStatus: %s", defConf.Title, defConf.Version, serverStatus)
	customMenus := PrepareMenuItems()
	trayInstance = systray.NewTrayManager(app, trayIcon, Tooltip, customMenus)

	if err := app.Run(); err != nil {
		log.Fatalf("❌ Linux system tray service failed: %v", err)
	}
}
