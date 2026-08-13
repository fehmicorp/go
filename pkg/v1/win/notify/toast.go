package notify

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-toast/toast"
)

var (
	appID          string
	defaultIcon    string
	onClickHandler func(string, string)
)

func RegisterClickCallback(handler func(string, string)) {
	onClickHandler = handler
}

func RegisterNotification(name string, iconPath string) error {
	if name == "" {
		return errors.New("application name / AppID cannot be empty")
	}

	appID = name

	if iconPath != "" {
		absPath, err := filepath.Abs(iconPath)
		if err != nil {
			return fmt.Errorf("failed to resolve absolute path for icon: %w", err)
		}
		if _, err := os.Stat(absPath); os.IsNotExist(err) {
			return fmt.Errorf("icon file does not exist at path: %s", absPath)
		}
		defaultIcon = absPath
	}

	return nil
}

// PushNotification sends a native Windows toast notification.
// Pass actionArgs in a format like "tag:action" (e.g., "app:onboard").
func PushNotification(title string, message string, actionArgs string) error {
	if appID == "" {
		return errors.New("notification package not registered: call RegisterNotification first")
	}
	notification := toast.Notification{
		AppID:               appID,
		Title:               title,
		Message:             message,
		Icon:                defaultIcon,
		ActivationArguments: actionArgs,
	}
	err := notification.Push()
	if err != nil {
		return fmt.Errorf("failed to push notification: %w", err)
	}
	return nil
}

func HandleClickCheck() bool {
	if len(os.Args) > 1 {
		rawArg := os.Args[1]
		if onClickHandler != nil {
			var tag, action string
			parts := strings.SplitN(rawArg, ":", 2)
			if len(parts) == 2 {
				tag = parts[0]
				action = parts[1]
			} else {
				tag = "default"
				action = rawArg
			}
			onClickHandler(tag, action)
			return true
		}
	}
	return false
}
