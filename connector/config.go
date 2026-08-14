package main

import (
	"embed"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

//go:embed assets/*
var AssetsFS embed.FS

type AppConfig struct {
	Title   string
	Desc    string
	Version string
	Cloud   string
}

var defConf = AppConfig{
	Title:   "Fehmi Cloud Connector",
	Desc:    "Standalone System Tray application",
	Version: "1.0.0",
	Cloud:   "connector.internal.svr",
}
var Target = struct {
	OS   string
	Arch string
}{
	OS:   runtime.GOOS,
	Arch: runtime.GOARCH,
}

type Asset struct {
	Name string
	Path string
}

var Assets = []Asset{
	{
		Name: "icon",
		Path: "/icon.png",
	},
}

func GetAssets(tag string) ([]byte, error) {
	var targetPath string
	for _, asset := range Assets {
		if asset.Name == tag {
			targetPath = asset.Path
			break
		}
	}
	if targetPath == "" {
		return nil, fmt.Errorf("asset not found: %s", tag)
	}
	targetDir := Target.OS + "-" + Target.Arch
	appDir := os.Getenv("APPDIR")
	if appDir != "" {
		fullPath := filepath.Join(appDir, "assets", targetDir, filepath.Base(targetPath))
		return os.ReadFile(fullPath)
	}
	return AssetsFS.ReadFile(targetPath)
}
