package utils

import (
	"fmt"
	"os"
)

// ReplaceOldVersion replaces the old version with the new version.
func ReplaceOldVersion(newPath, targetPath string) error {
	// 检查旧版本是否存在
	if _, err := os.Stat(targetPath); !os.IsNotExist(err) {
		backupPath := fmt.Sprintf("%s_backup", targetPath)
		if err := os.Rename(targetPath, backupPath); err != nil {
			return fmt.Errorf("failed to backup old version: %v", err)
		}
		fmt.Printf("Old version backed up to: %s\n", backupPath)
	}

	// 替换为新版本
	if err := os.Rename(newPath, targetPath); err != nil {
		return fmt.Errorf("failed to replace with new version: %v", err)
	}

	fmt.Printf("New version installed at: %s\n", targetPath)
	return nil
}
