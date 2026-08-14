package systray

import (
	"github.com/wailsapp/wails/v3/pkg/application"
)

type TrayManager struct {
	App  *application.App
	Tray *application.SystemTray
}

type MenuItemConfig struct {
	Title   string
	OnClick func(ctx *application.Context)
}

// NewTrayManager initializes the system tray with icon data, tooltips, and custom menus
func NewTrayManager(app *application.App, iconData []byte, tooltip string, customMenus []MenuItemConfig) *TrayManager {
	systray := app.SystemTray.New()

	if len(iconData) > 0 {
		systray.SetIcon(iconData)
	}
	if tooltip != "" {
		systray.SetTooltip(tooltip)
	}

	tm := &TrayManager{
		App:  app,
		Tray: systray,
	}

	tm.Refresh(tooltip, customMenus)

	return tm
}

// Refresh dynamically updates the system tray tooltip and rebuilds the menu items
func (tm *TrayManager) Refresh(tooltip string, customMenus []MenuItemConfig) {
	if tm.Tray == nil {
		return
	}

	if tooltip != "" {
		tm.Tray.SetTooltip(tooltip)
	}

	menu := application.NewMenu()

	// Add user-defined custom menus
	for _, item := range customMenus {
		title := item.Title
		handler := item.OnClick
		menu.Add(title).OnClick(func(ctx *application.Context) {
			if handler != nil {
				handler(ctx)
			}
		})
	}

	menu.AddSeparator()

	menu.Add("Close").OnClick(func(ctx *application.Context) {
		tm.App.Quit()
	})

	tm.Tray.SetMenu(menu)
}
