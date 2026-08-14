package main

import (
	"fmt"
	"log"

	"github.com/fehmicorp/go/pkg/v1/win/notify"
	"github.com/fehmicorp/go/pkg/v1/win/systray"
	"github.com/wailsapp/wails/v3/pkg/application"
)

func Windows() {
	app := application.New(application.Options{
		Name:        defConf.Title,
		Description: defConf.Desc,
	})
	err := WindowsNotification(app)
	if err != nil {
		log.Fatalf("Notification module Error: %v", err)
	}
	_ = app
}

func WindowsNotification(app *application.App) error {
	iconPath := "assets/icon.png"
	notify.RegisterNotification(defConf.Title, iconPath)
	notify.RegisterAlert(defConf.Title)
	notify.RegisterClickCallback(func(tag string, action string) {
		switch tag {
		default:
			fmt.Printf("Unhandled notification action: %s\n", action)
		}
	})
	notify.HandleClickCheck()

	Tooltip := fmt.Sprintf("%s\nVersion: %s\nStatus: %s", defConf.Title, defConf.Version, serverStatus)
	customMenus := PrepareMenuItems()
	trayIcon, err := AssetsFS.ReadFile(iconPath)
	if err != nil {
		log.Fatalf("failed to read embedded asset %s: %w", iconPath, err)
	}
	trayInstance = systray.NewTrayManager(app, trayIcon, Tooltip, customMenus)

	if err := app.Run(); err != nil {
		log.Fatalf("❌ System tray service failed: %v", err)
	}
	return nil
}
