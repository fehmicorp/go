package main

import (
	"log"

	"github.com/fehmicorp/go/pkg/v1/win/notify"
	"github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
	switch Target.OS {
	case "windows":
		Windows()
	case "linux":
		Linux()
	default:
		log.Fatalf("Unsupported operating system: %s", Target.OS)
	}
}

func Windows() {
	app := application.New(application.Options{
		Name:        defConf.Title,
		Description: defConf.Desc,
	})
	Icon, _ := GetAssets("icon")
	notify.RegisterNotification(defConf.Title, Icon)
	notify.RegisterAlert(defConf.Title)

}

func Linux() {

}
