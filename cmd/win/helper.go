package main

import (
	"os"
	"path/filepath"
)

func GetAssetPath(relativePath string) string {
	exePath, err := os.Executable()
	if err != nil {
		return relativePath
	}
	return filepath.Join(filepath.Dir(exePath), relativePath)
}
