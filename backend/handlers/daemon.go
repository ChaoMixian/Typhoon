// File: handlers/daemon.go
package handlers

import (
	"io"
	"net/http"

	"Typhoon/daemon"

	"github.com/gin-gonic/gin"
)

// StartMihomoHandler
func StartMihomoHandler(c *gin.Context) {
	// 直接调用 daemon 的封装方法
	if err := daemon.StartMihomoFromConfig(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Mihomo started successfully"})
}

// StopMihomoHandler
func StopMihomoHandler(c *gin.Context) {
	if err := daemon.StopMihomo(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Mihomo stopped successfully"})
}

// RestartMihomoHandler
func RestartMihomoHandler(c *gin.Context) {
	// 调用 daemon 的封装方法
	if err := daemon.RestartMihomoFromConfig(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Mihomo restarted successfully"})
}

// GetMihomoOutputHandler
func GetMihomoOutputHandler(c *gin.Context) {
	// 这里暂时保留原逻辑，若后续需要真实流式输出再做修改
	c.Stream(func(w io.Writer) bool {
		// 模拟输出
		w.Write([]byte("Mihomo is running...\n"))
		return true
	})
}

func GetMihomoStatusHandler(c *gin.Context) {
	// 检查 Mihomo 进程是否在运行
	if daemon.IsMihomoRunning() {
		c.JSON(http.StatusOK, gin.H{"status": "Mihomo is running"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Mihomo is not running"})
}
