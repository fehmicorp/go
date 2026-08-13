package main

import (
	"fmt"
	"os/exec"
)

var AppRunning bool

type AppList struct {
	Name string
	Path string
}

var Apps = []AppList{
	{
		Name: "onboard",
		Path: "app.exe",
	},
}

func LaunchApp(appName string) {
	var targetPath string
	for _, app := range Apps {
		if app.Name == appName {
			targetPath = app.Path
			break
		}
	}

	if targetPath == "" {
		fmt.Printf("App not found in registry: %s\n", appName)
		return
	}

	if !AppRunning {
		AppRunning = true
		go func() {
			defer func() { AppRunning = false }()
			fmt.Printf("Launching Dashboard App (%s)...\n", appName)
			cmd := exec.Command("./app/" + targetPath)
			if err := cmd.Start(); err != nil {
				fmt.Printf("Error starting app: %v\n", err)
				return
			}
			_ = cmd.Wait()
		}()
	}
}
