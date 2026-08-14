package main

import (
	"embed"
	_ "embed"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

//go:embed assets/*
var AssetsFS embed.FS

type AppConfig struct {
	Title     string
	Desc      string
	Version   string
	Cloud     string
	Webportal string
}

var Cfg AppSettings

type AppSettings struct {
	AtStartup bool   `json:"at_startup"`
	Port      string `json:"port"`
	Host      string `json:"host"`
	AppDir    string `json:"appDir"`
}

var (
	settingsFile = filepath.Join("./data", "config.json")
	currentCfg   AppSettings
	cfgMutex     sync.Mutex
)

func LoadConfig() AppSettings {
	cfgMutex.Lock()
	defer cfgMutex.Unlock()

	// Default settings: atStartup = true by default
	currentCfg = AppSettings{
		AtStartup: true,
		Port:      "8080",
		Host:      "localhost",
		AppDir:    "app",
	}

	file, err := os.Open(settingsFile)
	if err != nil {
		// Fallback to environment variables if config.json is not created yet
		if envPort := os.Getenv("PORT"); envPort != "" {
			currentCfg.Port = envPort
		}
		if envHost := os.Getenv("HOST"); envHost != "" {
			currentCfg.Host = envHost
		}
		if envHost := os.Getenv("APPDIR"); envHost != "" {
			currentCfg.AppDir = envHost
		}
		SaveConfigLocked(currentCfg)
		return currentCfg
	}
	defer file.Close()

	_ = json.NewDecoder(file).Decode(&currentCfg)

	// Fallback check if JSON has empty fields
	if currentCfg.Port == "" {
		currentCfg.Port = "8080"
	}
	if currentCfg.Host == "" {
		currentCfg.Host = "localhost"
	}
	if currentCfg.AppDir == "" {
		currentCfg.AppDir = "app"
	}
	return currentCfg
}

// SaveConfig updates and writes settings to disk
func SaveConfig(cfg AppSettings) {
	cfgMutex.Lock()
	defer cfgMutex.Unlock()
	SaveConfigLocked(cfg)
}

func SaveConfigLocked(cfg AppSettings) {
	currentCfg = cfg
	_ = os.MkdirAll(filepath.Dir(settingsFile), 0755)

	file, err := os.Create(settingsFile)
	if err != nil {
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	_ = encoder.Encode(cfg)
}

func GetConfig() AppSettings {
	cfgMutex.Lock()
	defer cfgMutex.Unlock()
	return currentCfg
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
