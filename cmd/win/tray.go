package main

import (
	"log"

	"github.com/fehmicorp/go/pkg/v1/win/notify"
	"github.com/fehmicorp/go/pkg/v1/win/systray"
	"github.com/wailsapp/wails/v3/pkg/application"
)

type MenuOptions struct {
	Title string
	Type  string
	Tag   string
	Func  func()
}

var MenuList = []MenuOptions{
	{
		Title: "Dashboard",
		Type:  "app",
		Tag:   "onboard",
	},
	{
		Title: "Push Notification",
		Type:  "notify",
		Tag:   "toast",
	},
	{
		Title: "Ring Bell",
		Type:  "notify",
		Tag:   "bell",
	},
}

func PrepareMenuItems() []systray.MenuItemConfig {
	var customMenus []systray.MenuItemConfig
	for _, menu := range MenuList {
		currentMenu := menu
		itemConfig := systray.MenuItemConfig{
			Title: currentMenu.Title,
			OnClick: func(ctx *application.Context) {
				if currentMenu.Func != nil {
					currentMenu.Func()
					return
				}
				switch currentMenu.Type {
				case "app":
					LaunchApp(currentMenu.Tag)
				case "notify":
					NotifyList(currentMenu.Tag)
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
