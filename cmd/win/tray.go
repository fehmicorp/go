package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"github.com/fehmicorp/go/pkg/v1/srvc"
	"github.com/fehmicorp/go/pkg/v1/win/notify"
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

func PrepareMenuItems() []systray.MenuItemConfig {
	menuList := []MenuOptions{
		{
			Title: "Dashboard",
			Type:  "web",
			Tag:   "http://localhost:8080",
			Func: func(ctx *application.Context) {
				OpenBrowser("http://localhost:8080")
			},
		},
		{
			Title:   "Start Services",
			Type:    "action",
			Tag:     "start",
			Dynamic: true,
			Func: func(ctx *application.Context) {
				fmt.Println("Tray action -> Starting services...")
				runServiceActionCommand("start")
			},
		},
		{
			Title: "Logs",
			Type:  "file",
			Tag:   "logs",
			Func: func(ctx *application.Context) {
				openLogFile()
			},
		},
		{
			Title: "Settings",
			Type:  "action",
			Tag:   "settings",
			Func: func(ctx *application.Context) {
				fmt.Println("Tray action -> Opening Settings")
			},
		},
	}
	var customMenus []systray.MenuItemConfig
	for _, menu := range menuList {
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

func NotifyList(tag string) {
	switch tag {
	case "toast":
		err := notify.PushNotification("Task Completed", "Your background process finished successfully!", "https://google.com/")
		if err != nil {
			log.Printf("Failed to send notification: %v", err)
		}
	case "bell":
		err := notify.BeepNotify("Beep", "Hello", Conf.Icon)
		if err != nil {
			log.Printf("Failed to send Beep: %v", err)
		}
	}
}

func OpenBrowser(url string) {
	var err error
	switch runtime.GOOS {
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = exec.Command("xdg-open", url).Start()
	}
	if err != nil {
		log.Printf("Failed to open browser: %v", err)
	}
}

func runServiceActionCommand(action string) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	baseURL := fmt.Sprintf("http://localhost:%s", port)
	switch action {
	case "start":
		resp, err := http.Get(fmt.Sprintf("%s/health", baseURL))
		if err != nil {
			log.Printf("Failed to reach server (is it running?): %v", err)
			return
		}
		defer resp.Body.Close()
		fmt.Printf("Service status check (/health) -> Status Code: %d\n", resp.StatusCode)
	case "stop":
		resp, err := http.Post(fmt.Sprintf("%s/stop", baseURL), "application/json", nil)
		if err != nil {
			hostOS, hostArch := runtime.GOOS, runtime.GOARCH
			log.Printf("Executing fallback service stop via srvc package...")
			list, fetchErr := srvc.FetchFilteredServices("", "", hostOS, hostArch)
			if fetchErr != nil {
				log.Printf("Failed to fetch services for stopping: %v", fetchErr)
				return
			}
			for i := range list {
				if stopErr := list[i].Stop(); stopErr != nil {
					log.Printf("Failed to stop service %s: %v", list[i].Name, stopErr)
				} else {
					fmt.Printf("Successfully stopped service: %s\n", list[i].Name)
				}
			}
			return
		}
		defer resp.Body.Close()
		fmt.Println("Server stop signal sent successfully.")
	default:
		fmt.Println("Initilizing")
	}

}

func openLogFile() {
	var err error
	if runtime.GOOS == "windows" {
		err = exec.Command("notepad.exe", "connector.log").Start()
	} else {
		err = exec.Command("xdg-open", "connector.log").Start()
	}
	if err != nil {
		log.Printf("Could not open log file: %v", err)
	}
}
