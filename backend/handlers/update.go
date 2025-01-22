// update mihomo or xray core, geoip, geosite, Typhoon, etc.

package handlers

import (
	"Typhoon/config"
	"Typhoon/update"
	"fmt"
	"log"
	"net/http"
	"path"

	"github.com/gin-gonic/gin"
)

// UpdateMihomoHandler handles the Mihomo update process
func UpdateMihomoHandler(c *gin.Context) {
	cfg := config.GetConfig(config.ConfigFilePath, true)

	// 从 URL 参数获取 downloadURL 和目标路径
	destPath := c.DefaultQuery("destPath", path.Join(cfg.Proxy.Mihomo.BinPath, "..", "mihomo_new"))
	finalPath := c.DefaultQuery("finalPath", path.Join(cfg.Proxy.Mihomo.BinPath, "..", "mihomo"))

	// 设置 SSE 头
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")

	// 更新 Mihomo 核心
	// 调用更新流程
	err := update.UpdateMihomo("MetaCubeX", "mihomo", destPath, finalPath, func(downloaded, total int64) {
		progress := float64(downloaded) / float64(total) * 100
		fmt.Printf("Progress: %.2f%%\n", progress)
	})
	if err != nil {
		log.Fatalf("Update failed: %v", err)
	}

	// Reload the configuration
	if config.GetConfig(config.ConfigFilePath, true) == nil {
		log.Fatalf("Failed to reload the configuration")
	}

	// 通知前端完成
	fmt.Fprintf(c.Writer, "data: {\"status\": \"completed\"}\n\n")
	c.Writer.(http.Flusher).Flush()
}

// UpdateSubscriptionsHandler handles the subscription update process
func UpdateSubscriptionsHandler(c *gin.Context) {
	// 更新订阅
	if err := update.UpdateSubscriptions(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "Subscriptions updated successfully"})
}
