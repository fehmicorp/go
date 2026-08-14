package main

import (
	"embed"
	_ "embed"
	"runtime"
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

var Status string

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
