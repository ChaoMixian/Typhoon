package utils

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
)

// DecompressGz decompresses a .gz file to the specified destination directory.
func DecompressGz(sourceFile, destFile string) error {
	// 打开 .gz 文件
	file, err := os.Open(sourceFile)
	if err != nil {
		return fmt.Errorf("failed to open source file: %v", err)
	}
	defer file.Close()

	// 创建 Gzip Reader
	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %v", err)
	}
	defer gzReader.Close()

	// 打开目标文件以写入
	outFile, err := os.Create(destFile)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %v", err)
	}
	defer outFile.Close()

	// 解压文件
	_, err = io.Copy(outFile, gzReader)
	if err != nil {
		return fmt.Errorf("failed to decompress file: %v", err)
	}

	return nil
}
