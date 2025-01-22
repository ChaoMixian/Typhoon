package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// Download downloads a file from the given URL and writes it to the specified destination.
func Download(url, destPath string) error {
	// 创建 HTTP 请求
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to initiate download: %v", err)
	}
	defer resp.Body.Close()

	// 检查 HTTP 响应状态
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download file: %v", resp.Status)
	}

	// 打开文件以写入
	destFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %v", err)
	}
	defer destFile.Close()

	// 将数据写入目标文件
	_, err = io.Copy(destFile, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to download file: %v", err)
	}

	return nil
}

// DownloadWithProgress downloads a file from the given URL and writes it to the specified destination.
// It reports progress via a callback function.
func DownloadWithProgress(url, destPath string, progressCallback func(downloaded, total int64)) error {
	// 确保目标文件夹存在
	dir := filepath.Dir(destPath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	// 创建 HTTP 请求
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to initiate download: %v", err)
	}
	defer resp.Body.Close()

	// 检查 HTTP 响应状态
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download file: %v", resp.Status)
	}

	// 获取文件大小
	totalSize := resp.ContentLength

	// 打开文件以写入
	destFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %v", err)
	}
	defer destFile.Close()

	// 创建一个带进度跟踪的 Reader
	progressReader := &progressReader{
		Reader: resp.Body,
		Total:  totalSize,
		Callback: func(downloaded int64) {
			if progressCallback != nil {
				progressCallback(downloaded, totalSize)
			}
		},
	}

	// 将数据写入目标文件
	_, err = io.Copy(destFile, progressReader)
	if err != nil {
		return fmt.Errorf("failed to download file: %v", err)
	}

	return nil
}

// progressReader tracks the progress of a download.
type progressReader struct {
	Reader     io.Reader
	Total      int64
	Downloaded int64
	Callback   func(downloaded int64)
}

// Read updates the progress and calls the callback function.
func (p *progressReader) Read(b []byte) (int, error) {
	n, err := p.Reader.Read(b)
	p.Downloaded += int64(n)
	if p.Callback != nil {
		p.Callback(p.Downloaded)
	}
	return n, err
}
