package notify

import "github.com/gen2brain/beeep"

func RegisterAlert(Appname string) {
	beeep.AppName = Appname
}

func BeepAlert(Title string, Msg string, IconPath string) error {
	return beeep.Alert(Title, Msg, IconPath)
}

func BeepNotify(Title string, Msg string, IconPath string) error {
	return beeep.Notify(Title, Msg, IconPath)
}

func SimpleAlert(Title string, Msg string) error {
	return beeep.Notify(Title, Msg, "")
}
