package handlers

import (
	"Typhoon/config"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// convertToStringSlice converts an interface{} to a []string slice
func convertToStringSlice(value interface{}) []string {
	values, ok := value.([]interface{})
	if !ok {
		return nil
	}

	result := make([]string, len(values))
	for i, v := range values {
		result[i] = v.(string)
	}
	return result
}

// convertToSubscriptionSlice converts an interface{} to a subscription slice
func convertToSubscriptionSlice(value interface{}) []struct {
	Name string `json:"name"`
	URL  string `json:"url"`
	Path string `json:"path"`
} {
	values, ok := value.([]interface{})
	if !ok {
		return nil
	}

	result := make([]struct {
		Name string `json:"name"`
		URL  string `json:"url"`
		Path string `json:"path"`
	}, len(values))

	for i, v := range values {
		sub := v.(map[string]interface{})
		path, ok := sub["path"].(string)
		if !ok {
			path = "" // 如果没有提供 path 字段，则设置为空字符串
		}
		result[i] = struct {
			Name string `json:"name"`
			URL  string `json:"url"`
			Path string `json:"path"`
		}{
			Name: sub["name"].(string),
			URL:  sub["url"].(string),
			Path: path,
		}
	}
	return result
}

// ReloadConfigHandler handles API requests to reload configuration
func ReloadConfigHandler(c *gin.Context) {
	// 重新加载配置
	if config.GetConfig(config.ConfigFilePath, true) == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reload the configuration"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Configuration reloaded successfully"})
}

// UpdateConfigHandler handles API requests to update configuration
func UpdateConfigHandler(c *gin.Context) {
	// 获取配置文件路径
	filePath := c.Query("filePath")
	if filePath == "" {
		filePath = config.ConfigFilePath // 默认配置文件路径
	}

	// 获取要更新的字段和值
	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// 更新配置
	err := config.UpdateConfig(filePath, func(cfg *config.Config) error {
		for key, value := range updates {
			switch key {
			// 代理配置
			case "proxy.currentCore":
				cfg.Proxy.CurrentCore = value.(string)
			case "proxy.mihomo.listenPort":
				cfg.Proxy.Mihomo.ListenPort = int(value.(float64))
			case "proxy.mihomo.binPath":
				cfg.Proxy.Mihomo.BinPath = value.(string)
			case "proxy.mihomo.currentConfig":
				cfg.Proxy.Mihomo.CurrentConfig = value.(string)
			case "proxy.mihomo.controllerAddress":
				cfg.Proxy.Mihomo.ControllerAddress = value.(string)
			case "proxy.mihomo.tun.enabled":
				cfg.Proxy.Mihomo.TUN.Enabled = value.(bool)
			case "proxy.mihomo.tun.stack":
				cfg.Proxy.Mihomo.TUN.Stack = value.(string)
			case "proxy.mihomo.tun.dnsHijaking":
				cfg.Proxy.Mihomo.TUN.DNSHijaking = value.(string)

			// DNS 配置
			case "proxy.dns.enabled":
				cfg.Proxy.DNS.Enabled = value.(bool)
			case "proxy.dns.listen":
				cfg.Proxy.DNS.Listen = value.(string)
			case "proxy.dns.upstreamDNS":
				cfg.Proxy.DNS.UpstreamDNS = convertToStringSlice(value)
			case "proxy.dns.fallbackDNS":
				cfg.Proxy.DNS.FallbackDNS = convertToStringSlice(value)
			case "proxy.dns.dnsHijaking.enabled":
				cfg.Proxy.DNS.DNSHijaking.Enabled = value.(bool)
			case "proxy.dns.enhancedMode":
				cfg.Proxy.DNS.EnhancedMode = value.(string)
			case "proxy.dns.fakeIPFilter":
				cfg.Proxy.DNS.FakeIPFilter = convertToStringSlice(value)

			// 日志配置
			case "logging.level":
				cfg.Logging.Level = value.(string)
			case "logging.file":
				cfg.Logging.File = value.(string)

			// 订阅配置
			case "subscriptionManage.enabled":
				cfg.SubscriptionManage.Enabled = value.(bool)
			case "subscriptionManage.intervalSeconds":
				cfg.SubscriptionManage.IntervalSeconds = int(value.(float64))
			case "subscriptionManage.subscriptions":
				cfg.SubscriptionManage.Subscriptions = convertToSubscriptionSlice(value)

			// API 配置
			case "api.listenPort":
				cfg.API.ListenPort = int(value.(float64))
			case "api.token":
				cfg.API.Token = value.(string)

			default:
				return fmt.Errorf("unsupported config key: %s", key)
			}
		}
		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Configuration updated successfully"})
}

// Todo: 抽象订阅链接操作，增删改查

/*
curl -X PATCH http://localhost:8080/api/v1/updateConfig \
-H "Content-Type: application/json" \
-d '{
    "proxy.currentCore": "xray",
    "proxy.mihomo.listenPort": 8888
}'

curl -X PATCH http://localhost:8080/api/v1/updateConfig \
-H "Content-Type: application/json" \
-d '{
    "proxy.dns.enabled": true,
    "proxy.dns.listen": "127.0.0.1:5353",
    "proxy.dns.upstreamDNS": ["1.1.1.1", "8.8.8.8"],
    "proxy.dns.fakeIPFilter": ["*.example.com", "localhost"]
}'


curl -X PATCH http://localhost:8080/api/v1/updateConfig \
-H "Content-Type: application/json" \
-d '{
    "subscriptionManage.enabled": true,
    "subscriptionManage.subscriptions": [
        {"name": "Sub1", "url": "https://example.com/sub1"},
        {"name": "Sub2", "url": "https://example.com/sub2"}
    ]
}'

*/
