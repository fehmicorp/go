package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/fehmicorp/go/pkg/v1/srvc"
)

type App struct {
	ctx            context.Context
	supabaseURL    string
	supabaseKey    string
	cachedServices []srvc.Services
}

func NewApp() *App {
	key := os.Getenv("SUPABASE_SECRET_KEY")
	if key == "" {
		key = os.Getenv("SUPABASE_PUBLISHABLE_KEY")
	}
	if key == "" {
		key = os.Getenv("SUPABASE_ANON_KEY")
	}
	return &App{
		supabaseURL: os.Getenv("SUPABASE_URL"),
		supabaseKey: key,
	}
}

func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	targetOS, targetArch := getTarget()
	list, err := srvc.FetchFilteredServices(a.supabaseURL, a.supabaseKey, targetOS, targetArch)
	if err != nil {
		log.Printf("Warning: Failed to fetch services from Supabase on startup: %v", err)
		a.cachedServices = []srvc.Services{}
		return
	}
	a.cachedServices = list
}

func getTarget() (string, string) {
	return runtime.GOOS, runtime.GOARCH
}

func (a *App) GetServiceList() []srvc.Services {
	targetOS, targetArch := getTarget()
	list, err := srvc.FetchFilteredServices(a.supabaseURL, a.supabaseKey, targetOS, targetArch)
	if err != nil {
		log.Printf("Failed to fetch filtered services: %v", err)
		return a.cachedServices
	}

	// Update live system statuses using the enhanced RefreshStatus logic
	for i := range list {
		_ = list[i].RefreshStatus()
	}

	return list
}

// CheckServiceVersion inspects the host system for a specific service's installation status and version updates
func (a *App) CheckServiceVersion(name string) (srvc.VersionCheckResult, error) {
	var targetService *srvc.Services
	targetOS, targetArch := getTarget()

	list, err := srvc.FetchFilteredServices(a.supabaseURL, a.supabaseKey, targetOS, targetArch)
	if err != nil {
		list = a.cachedServices
	}

	for i := range list {
		if list[i].Name == name {
			targetService = &list[i]
			break
		}
	}

	if targetService == nil {
		return srvc.VersionCheckResult{}, fmt.Errorf("service with name '%s' not found", name)
	}

	return targetService.CheckLocalVersion()
}

func (a *App) ServiceAction(name string, action string) error {
	var targetService *srvc.Services
	targetOS, targetArch := getTarget()
	list, err := srvc.FetchFilteredServices(a.supabaseURL, a.supabaseKey, targetOS, targetArch)
	if err != nil {
		list = a.cachedServices
	}
	for i := range list {
		if list[i].Name == name {
			targetService = &list[i]
			break
		}
	}
	if targetService == nil {
		return fmt.Errorf("service with name '%s' not found", name)
	}
	switch action {
	case "start":
		return targetService.Start()
	case "stop":
		return targetService.Stop()
	case "restart":
		return targetService.Restart()
	case "install":
		execPath := fmt.Sprintf("%s/%s", targetService.Package.WorkDir, targetService.Package.Name)
		return targetService.Install(execPath)
	case "uninstall":
		return targetService.Uninstall()
	default:
		return fmt.Errorf("unknown action type: %s", action)
	}
}
