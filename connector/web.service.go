package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/fehmicorp/go/pkg/v1/logger"
)

func StartServerLocked() {
	if isServerRunning {
		return
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	host := os.Getenv("HOST")
	if host == "" {
		host = "localhost"
	}

	addr := fmt.Sprintf("%s:%s", host, port)
	httpServer = &http.Server{
		Addr:    addr,
		Handler: nil,
	}

	isServerRunning = true
	serverStatus = "Running"

	go func() {
		fmt.Printf("Starting web server on http://%s/\n", addr)
		logger.Info("SERVER", fmt.Sprintf("Web server started on http://%s/", addr), nil)

		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("SERVER", "Web server stopped with error", map[string]interface{}{"error": err.Error()})
		}
		serverMutex.Lock()
		isServerRunning = false
		serverStatus = "Stopped"
		serverMutex.Unlock()
		updateTrayState()
	}()

	for i := range TrayMenu {
		if TrayMenu[i].Tag == "start" {
			TrayMenu[i].Title = "Stop"
			break
		}
	}
	fmt.Println("Tray action -> Services started.")
}

func StopServerLocked() {
	if !isServerRunning {
		return
	}

	if httpServer != nil {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		_ = httpServer.Shutdown(shutdownCtx)
		httpServer = nil
	}

	isServerRunning = false
	serverStatus = "Stopped"
	logger.Info("SERVER", "Web server stopped by user action", nil)

	for i := range TrayMenu {
		if TrayMenu[i].Tag == "start" {
			TrayMenu[i].Title = "Start"
			break
		}
	}
	fmt.Println("Tray action -> Services stopped.")
}
