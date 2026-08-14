package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/fehmicorp/go/pkg/v1/win/systray"
	"github.com/wailsapp/wails/v3/pkg/application"
)

type MenuOptions struct {
	Title   string
	Type    string
	Tag     string
	Dynamic bool
	Func    func(*application.Context)
}

var (
	serverMutex     sync.Mutex
	httpServer      *http.Server
	isServerRunning bool
	serverStatus    = "Stopped"
	trayInstance    *systray.TrayManager
)

var TrayMenu []MenuOptions

func initTrayMenu() {
	serverMutex.Lock()
	serverStatusTitle := "Start"
	if isServerRunning {
		serverStatusTitle = "Stop"
	}
	serverMutex.Unlock()
	TrayMenu = []MenuOptions{
		{
			Title: "Dashboard",
			Type:  "web",
			Tag:   "http://localhost:8080",
			Func: func(ctx *application.Context) {
				serverMutex.Lock()
				if !isServerRunning {
					fmt.Println("Server is not running. Starting service automatically...")
					StartServerLocked()
				}
				portalURL := defConf.Webportal
				serverMutex.Unlock()

				time.Sleep(300 * time.Millisecond)

				fmt.Printf("Opening Dashboard -> %s\n", portalURL)
				OpenBrowser(portalURL)
			},
		},
		{
			Title:   serverStatusTitle,
			Type:    "action",
			Tag:     "start",
			Dynamic: true,
			Func: func(ctx *application.Context) {
				serverMutex.Lock()
				if !isServerRunning {
					StartServerLocked()
				} else {
					StopServerLocked()
				}
				serverMutex.Unlock()
				updateTrayState()
			},
		},
		{
			Title:   getAutoStartMenuTitle(Cfg.AtStartup),
			Type:    "action",
			Tag:     "toggle_autostart",
			Dynamic: true,
			Func: func(ctx *application.Context) {
				cfgMutex.Lock()
				Cfg.AtStartup = !Cfg.AtStartup
				SaveConfigLocked(Cfg)
				updatedAtStartup := Cfg.AtStartup
				cfgMutex.Unlock()
				fmt.Printf("Auto-start on boot updated to: %v\n", updatedAtStartup)
				initTrayMenu()
				updateTrayState()
			},
		},
		{
			Title: "Logs",
			Type:  "file",
			Tag:   "logs",
			Func: func(ctx *application.Context) {
				fmt.Println("Exporting and opening logs...")
				openLogsViewer()
			},
		},
	}
}

func getAutoStartMenuTitle(atStartup bool) string {
	if atStartup {
		return "✓ Start at System Startup"
	}
	return "Start at System Startup"
}

func updateTrayState() {
	if trayInstance != nil {
		tooltip := fmt.Sprintf("%s\nVersion: %s\nStatus: %s", defConf.Title, defConf.Version, serverStatus)
		trayInstance.Refresh(tooltip, PrepareMenuItems())
	}
}

func PrepareMenuItems() []systray.MenuItemConfig {
	var customMenus []systray.MenuItemConfig
	for _, menu := range TrayMenu {
		currentMenu := menu
		itemConfig := systray.MenuItemConfig{
			Title: currentMenu.Title,
			OnClick: func(ctx *application.Context) {
				if currentMenu.Func != nil {
					currentMenu.Func(ctx)
					return
				}
			},
		}
		customMenus = append(customMenus, itemConfig)
	}
	return customMenus
}
