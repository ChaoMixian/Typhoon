package utils

import (
	"log"
	"os"
	"path/filepath"
)

// GetExecutableDir returns the directory of the executable file
func GetExecutableDir() string {
	executablePath, err := os.Executable()
	if err != nil {
		log.Fatalf("Failed to get executable path: %v", err)
	}
	execDir := filepath.Dir(executablePath)
	return execDir
}
