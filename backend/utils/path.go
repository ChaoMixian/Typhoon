package utils

import (
	"log"
	"os"
	"path/filepath"
)

// GetExecutableDir returns the directory of the executable file
func GetExecutableDir() string {
	executablePath, err := os.Executable()
	// fmt.Printf("executablePath: %v\n", executablePath)
	if err != nil {
		log.Printf("Failed to get executable path: %v", err)
	}
	execDir := filepath.Dir(executablePath)
	return execDir
}

// GetMihomoConfigPath returns the path of the Mihomo configuration file
func GetMihomoConfigPath(configName string) (string, string) {
	execDir := GetExecutableDir()
	return filepath.Join(execDir, "mihomo", "config", configName, "config.yaml"),
		filepath.Join(execDir, "mihomo", "config", configName, "runtime.yaml")
}
