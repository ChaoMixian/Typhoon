package handlers

import (
	"Typhoon/config"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

// ReverseProxyHandler handles the reverse proxy for Mihomo API
func ReverseProxyHandler(c *gin.Context) {
	// 定义 Mihomo API 地址
	cfg, err := config.GetConfig(config.ConfigFilePath, true)
	target := "http://" + cfg.Proxy.Mihomo.ControllerAddress
	remote, err := url.Parse(target)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse target URL"})
		return
	}

	// 创建反向代理
	proxy := httputil.NewSingleHostReverseProxy(remote)
	// 重写路径，移除 /api/v1/mihomo 并保留后续路径
	proxy.Director = func(req *http.Request) {
		req.URL.Scheme = remote.Scheme
		req.URL.Host = remote.Host
		req.URL.Path = c.Param("path") // 获取 /api/v1/mihomo/*path 的 *path 部分
		req.Host = remote.Host
	}
	proxy.ServeHTTP(c.Writer, c.Request)
}
