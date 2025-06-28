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
	cfg, cfg_err := config.GetConfig(config.ConfigFilePath, true)
	if cfg_err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve current configuration: " + cfg_err.Error()})
		return
	}
	// log.Println(cfg.Proxy.Mihomo.BinPath)
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
		log.Printf("Update failed: %v", err)
		fmt.Fprintf(c.Writer, "data: {\"status\": \"failed\", \"error\": \"%s\"}\n\n", err.Error())
		c.Writer.(http.Flusher).Flush()
	}

	// Reload the configuration
	_, reloadErr := config.GetConfig(config.ConfigFilePath, true)
	if reloadErr != nil {
		log.Fatalf("Failed to reload the configuration: %v", reloadErr)
		fmt.Fprintf(c.Writer, "data: {\"status\": \"failed\", \"error\": \"Failed to reload the configuration\"}\n\n")
		c.Writer.(http.Flusher).Flush()
	}

	// 通知前端完成
	fmt.Fprintf(c.Writer, "data: {\"status\": \"completed\"}\n\n")
	c.Writer.(http.Flusher).Flush()
}

// UpdateSubscriptionsHandler handles the subscription update process
func UpdateSubscriptionsHandler(c *gin.Context) {
	// 更新订阅
	results, err := update.UpdateSubscriptions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "配置保存失败",
			"details": err.Error(),
			"results": results,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "Subscriptions updated successfully",
		"results": results,
	})
}
