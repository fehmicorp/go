package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

// openLogsViewer reads sqlite records, formats them into a text report, and opens it in Notepad/Editor
func openLogsViewer() {
	dbPath := filepath.Join("./data", "logs.db")
	db, err := sql.Open("sqlite3", dbPath+"?mode=ro") // Open database in read-only mode to prevent locks
	if err != nil {
		log.Printf("Failed to open logs database for reading: %v", err)
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT level, category, message, metadata, created_at FROM application_logs ORDER BY id DESC LIMIT 500")
	if err != nil {
		log.Printf("Failed to query logs: %v", err)
		return
	}
	defer rows.Close()

	// Create a temp readable text file
	tmpFile, err := os.CreateTemp("", "fehmi_logs_*.txt")
	if err != nil {
		log.Printf("Failed to create temporary log text file: %v", err)
		return
	}
	defer tmpFile.Close()

	// Write header info
	_, _ = tmpFile.WriteString("=================================================================\n")
	_, _ = tmpFile.WriteString("               FEHMI CLOUD CONNECTOR - APPLICATION LOGS           \n")
	_, _ = tmpFile.WriteString("=================================================================\n\n")

	count := 0
	for rows.Next() {
		var level, category, message, metadata, createdAt string
		if err := rows.Scan(&level, &category, &message, &metadata, &createdAt); err == nil {
			count++
			line := fmt.Sprintf("[%s] [%s] (%s) - %s\nMetadata: %s\n-----------------------------------------------------------------\n",
				createdAt, level, category, message, metadata)
			_, _ = tmpFile.WriteString(line)
		}
	}

	if count == 0 {
		_, _ = tmpFile.WriteString("No logs found in database yet.\n")
	}

	// Open the generated text file using the OS text editor (Notepad on Windows)
	var cmd *exec.Cmd
	switch Target.OS {
	case "windows":
		cmd = exec.Command("notepad.exe", tmpFile.Name())
	case "darwin":
		cmd = exec.Command("open", tmpFile.Name())
	default:
		cmd = exec.Command("xdg-open", tmpFile.Name())
	}

	if err := cmd.Start(); err != nil {
		log.Printf("Failed to open text editor for logs: %v", err)
	}
}

func OpenBrowser(url string) {
	var err error
	switch Target.OS {
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = exec.Command("xdg-open", url).Start()
	}
	if err != nil {
		log.Printf("Failed to open browser: %v", err)
	}
}
