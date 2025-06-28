package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// GetExecutableDir returns the directory of the executable file
func GetExecutableDir() (string, error) {
	executablePath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get executable path: %w", err)
	}
	execDir := filepath.Dir(executablePath)
	return execDir, nil
}

// GetMihomoConfigPath returns the path of the Mihomo configuration file.
// It logs a fatal error if GetExecutableDir fails, as these paths are critical.
func GetMihomoConfigPath(configName string) (string, string) {
	execDir, err := GetExecutableDir()
	if err != nil {
		log.Fatalf("Critical error: Could not determine executable directory for Mihomo config paths: %v", err)
	}
	return filepath.Join(execDir, "mihomo", "config", configName, "config.yaml"),
		filepath.Join(execDir, "mihomo", "config", configName, "runtime.yaml")
}

func IsSubPath(base, sub string) bool {
	// Check if sub is a subdirectory of base
	rel, err := filepath.Rel(base, sub)
	if err != nil {
		return false // If we can't determine the relative path, return false
	}
	return !filepath.IsAbs(rel) // If rel is not absolute, sub is a subpath of base
}
