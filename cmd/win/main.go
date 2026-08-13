package main

import (
	_ "embed"
	"fmt"
	"log"

	"github.com/fehmicorp/go/pkg/v1/win/notify"
	"github.com/fehmicorp/go/pkg/v1/win/systray"
	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed assets/icon.png
var trayIcon []byte

type AppConfig struct {
	AppName     string
	Description string
	Icon        string
	Version     string
	Domain      string
}

var Conf = AppConfig{
	AppName:     "Fehmi Cloud Connector",
	Description: "Standalone Windows System Tray application",
	Icon:        "assets/icon.png",
	Version:     "1.0.0",
	Domain:      "fehmicorp.in",
}

func main() {
	app := application.New(application.Options{
		Name:        Conf.AppName,
		Description: Conf.Description,
	})

	err := notify.RegisterNotification(Conf.AppName, Conf.Icon)
	notify.RegisterAlert(Conf.AppName)
	if err != nil {
		log.Fatalf("Registration failed: %v", err)
	}

	notify.RegisterClickCallback(func(tag string, action string) {
		fmt.Printf("Notification clicked -> Tag: %s, Action: %s\n", tag, action)
		switch tag {
		case "app":
			LaunchApp(action)
		case "notify":
			switch action {
			case "toast":
				fmt.Println("Toast notification interaction handled.")
			case "bell":
				fmt.Println("Bell notification interaction handled.")
			default:
				fmt.Printf("Unhandled notification action: %s\n", action)
			}
		default:
			fmt.Printf("Unhandled notification action: %s\n", action)
		}
	})

	notify.HandleClickCheck()

	Tooltip := fmt.Sprintf("%s\nVersion: %s\n%s", Conf.AppName, Conf.Version, Conf.Domain)
	customMenus := PrepareMenuItems()

	systray.NewTrayManager(app, trayIcon, Tooltip, customMenus)

	if err := app.Run(); err != nil {
		log.Fatalf("❌ System tray service failed: %v", err)
	}
}
